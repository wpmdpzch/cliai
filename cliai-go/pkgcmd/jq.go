package pkgcmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ExecJq 执行 jq
func ExecJq(args []string) error {
	var buf strings.Builder
	err := ExecJqToWriter(args, &buf)
	if err != nil {
		return err
	}
	fmt.Print(buf.String())
	return nil
}

// ExecJqToWriter 执行 jq 并写入指定 writer
func ExecJqToWriter(args []string, w io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: jq [options] <filter> <file>")
	}

	// 简单参数解析
	var filter string
	var filename string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if filter == "" {
			filter = arg
		} else {
			filename = arg
		}
	}

	if filter == "" {
		filter = "."
	}

	// 读取输入
	var data []byte
	var err error

	if filename == "" || filename == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(filename)
	}

	if err != nil {
		return err
	}

	// 解析 JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("JSON 解析错误: %v", err)
	}

	// 应用过滤器
	result, err := applyFilter(jsonData, filter)
	if err != nil {
		return err
	}

	// 输出结果
	if result == nil {
		return nil
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(output))

	return nil
}

// applyFilter 应用 jq 风格的过滤器
func applyFilter(data interface{}, filter string) (interface{}, error) {
	filter = strings.TrimPrefix(filter, ".")

	if filter == "" {
		return data, nil
	}

	parts := strings.Split(filter, ".")
	current := data

	for _, part := range parts {
		if part == "" {
			continue
		}

		// 处理数组索引，如 [0] 或 [0:2]
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			idxStr := part[1 : len(part)-1]
			
			// 支持 [n] 格式
			if idx, err := strconv.Atoi(idxStr); err == nil {
				switch v := current.(type) {
				case []interface{}:
					if idx < 0 {
						idx = len(v) + idx
					}
					if idx < 0 || idx >= len(v) {
						return nil, nil
					}
					current = v[idx]
				default:
					return nil, fmt.Errorf("无法对 %T 应用数组索引", current)
				}
			} else if idxStr == "" {
				// [\] 表示全部元素
				current = current
			} else {
				return nil, fmt.Errorf("不支持的数组索引格式: %s", idxStr)
			}
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[part]; ok {
				current = val
			} else {
				return nil, nil
			}
		case []interface{}:
			// 尝试将 part 作为数字索引
			if idx, err := strconv.Atoi(part); err == nil {
				if idx < 0 {
					idx = len(v) + idx
				}
				if idx < 0 || idx >= len(v) {
					return nil, nil
				}
				current = v[idx]
			} else {
				return nil, fmt.Errorf("无法在数组上使用字段名 '%s'", part)
			}
		default:
			return nil, fmt.Errorf("无法在 %T 上应用字段 '%s'", current, part)
		}
	}

	return current, nil
}
