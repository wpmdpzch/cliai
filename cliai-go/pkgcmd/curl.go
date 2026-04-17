package pkgcmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ExecCurl 执行 curl
func ExecCurl(args []string) error {
	var buf strings.Builder
	err := ExecCurlToWriter(args, &buf)
	if err != nil {
		return err
	}
	fmt.Print(buf.String())
	return nil
}

// ExecCurlToWriter 执行 curl 并写入指定 writer
func ExecCurlToWriter(args []string, w io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: curl [options] <url>")
	}

	// 简单参数解析
	method := "GET"
	var url string
	headers := make(map[string]string)
	showHeader := false
	outputFile := ""
	var data string

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
		case "-d", "--data", "--data-raw":
			if i+1 < len(args) {
				data = args[i+1]
				method = "POST"
				if _, ok := headers["Content-Type"]; !ok {
					headers["Content-Type"] = "application/x-www-form-urlencoded"
				}
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
	var body io.Reader
	if data != "" {
		body = strings.NewReader(data)
	}
	req, err := http.NewRequest(method, url, body)
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
		fmt.Fprintf(w, "HTTP/%d %s\n", resp.StatusCode, resp.Status)
		for k, v := range resp.Header {
			fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Fprintln(w)
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
		_, err = io.Copy(w, resp.Body)
	}

	return err
}
