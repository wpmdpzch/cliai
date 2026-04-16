package pkgcmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// ExecJq Go 原生 jq 实现（简化版）
func ExecJq(args []string) error {
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
	fmt.Println(string(output))

	return nil
}

// applyFilter 应用 jq 风格的过滤器
func applyFilter(data interface{}, filter string) (interface{}, error) {
	// 简化实现：支持 .field 和 .field.subfield
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

		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[part]; ok {
				current = val
			} else {
				return nil, nil
			}
		case []interface{}:
			// 支持 [0] 索引
				if part == "@uniq" {
					current = uniq(v)
				} else {
					return nil, fmt.Errorf("数组索引需要 [n] 格式")
				}
		default:
			return nil, fmt.Errorf("无法在 %T 上应用字段 '%s'", current, part)
		}
	}

	return current, nil
}

// uniq 去重
func uniq(arr []interface{}) []interface{} {
	seen := make(map[string]bool)
	result := []interface{}{}
	for _, v := range arr {
		key := fmt.Sprintf("%v", v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return result
}
