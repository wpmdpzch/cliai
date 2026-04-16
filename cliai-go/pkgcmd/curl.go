package pkgcmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ExecCurl Go 原生 curl 实现
func ExecCurl(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: curl [options] <url>")
	}

	// 简单参数解析
	method := "GET"
	var url string
	headers := make(map[string]string)
	showHeader := false
	outputFile := ""

	i := 0
	for i < len(args) {
		arg := args[i]
		switch arg {
		case "-X", "--request":
			if i+1 < len(args) {
				method = args[i+1]
				i++
			}
		case "-H", "--header":
			if i+1 < len(args) {
				header := args[i+1]
				if idx := strings.Index(header, ":"); idx != -1 {
					headers[header[:idx]] = strings.TrimSpace(header[idx+1:])
				}
				i++
			}
		case "-i", "--include":
			showHeader = true
		case "-o", "--output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-d", "--data":
			if i+1 < len(args) {
				method = "POST"
				headers["Content-Type"] = "application/x-www-form-urlencoded"
				i++
			}
		default:
			if !strings.HasPrefix(arg, "-") {
				url = arg
			}
		}
		i++
	}

	if url == "" {
		return fmt.Errorf("错误: 请指定 URL")
	}

	// 添加协议前缀
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	// 创建请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 输出响应头
	if showHeader {
		fmt.Printf("HTTP/%d %s\n", resp.StatusCode, resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Println()
	}

	// 输出响应体
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
	} else {
		_, err = io.Copy(os.Stdout, resp.Body)
	}

	return err
}
