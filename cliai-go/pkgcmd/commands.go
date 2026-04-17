package pkgcmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Command 命令定义
type Command struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Example     string `json:"example"`
	Implemented string `json:"implemented"` // "go" or "system"
	Dangerous   bool   `json:"dangerous"`
	Args        []string `json:"args,omitempty"`
}

// commands 全局命令注册表
var commands = make(map[string]*Command)

func init() {
	registerP0Commands()
	registerP1Commands()
}

// registerP0Commands 注册 P0 命令
func registerP0Commands() {
	cmds := []Command{
		{Name: "curl", Category: "network", Description: "发送 HTTP/HTTPS 请求", Usage: "curl [options] <url>", Example: "curl https://example.com", Implemented: "go"},
		{Name: "jq", Category: "text", Description: "JSON 处理", Usage: "jq [options] <filter> <file>", Example: "jq '.data' file.json", Implemented: "go"},
		{Name: "grep", Category: "text", Description: "文本搜索", Usage: "grep [options] <pattern> <file>", Example: "grep 'error' log.txt", Implemented: "system"},
		{Name: "cat", Category: "text", Description: "查看文件内容", Usage: "cat <file>", Example: "cat readme.md", Implemented: "system"},
		{Name: "ls", Category: "file", Description: "列出目录", Usage: "ls [options] <dir>", Example: "ls -la", Implemented: "system"},
		{Name: "head", Category: "file", Description: "查看文件头部", Usage: "head [options] <file>", Example: "head -n 10 file.txt", Implemented: "system"},
		{Name: "tail", Category: "file", Description: "查看文件尾部", Usage: "tail [options] <file>", Example: "tail -n 10 file.txt", Implemented: "system"},
		{Name: "ps", Category: "system", Description: "查看进程", Usage: "ps [options]", Example: "ps aux", Implemented: "system"},
		{Name: "df", Category: "system", Description: "查看磁盘使用", Usage: "df [options]", Example: "df -h", Implemented: "system"},
		{Name: "free", Category: "system", Description: "查看内存使用", Usage: "free [options]", Example: "free -h", Implemented: "system"},
		{Name: "base64", Category: "encoding", Description: "Base64 编解码", Usage: "base64 [options] <file>", Example: "base64 file.txt", Implemented: "go"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}

// ExecResult 命令执行结果
type ExecResult struct {
	Output string
	Error  error
}

// Exists 检查命令是否存在
func Exists(name string) bool {
	_, ok := commands[name]
	return ok
}

// Get 获取命令
func Get(name string) *Command {
	return commands[name]
}

// List 返回所有命令
func List() []*Command {
	result := make([]*Command, 0, len(commands))
	for _, cmd := range commands {
		result = append(result, cmd)
	}
	return result
}

// ExecCommand 执行任意命令，返回结果
func ExecCommand(input string) *ExecResult {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return &ExecResult{Output: "", Error: nil}
	}

	name := parts[0]
	cmd, ok := commands[name]
	if !ok {
		return execSystemCmd(input)
	}

	if cmd.Implemented == "go" {
		return execGoCmd(cmd, parts[1:])
	}

	return execSystemCmd(input)
}

// execGoCmd 执行 Go 原生命令
func execGoCmd(cmd *Command, args []string) *ExecResult {
	var buf bytes.Buffer

	switch cmd.Name {
	case "base64":
		err := ExecBase64ToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	case "curl":
		err := ExecCurlToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	case "jq":
		err := ExecJqToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	default:
		return &ExecResult{Output: "", Error: fmt.Errorf("未实现的 Go 命令: %s", cmd.Name)}
	}
}

// execSystemCmd 执行系统命令
func execSystemCmd(input string) *ExecResult {
	args := strings.Fields(input)
	if len(args) == 0 {
		return &ExecResult{Output: "", Error: nil}
	}

	cmd := exec.Command(args[0], args[1:]...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	return &ExecResult{Output: buf.String(), Error: err}
}

// registerP1Commands 注册 P1 命令
func registerP1Commands() {
	cmds := []Command{
		{Name: "wget", Category: "network", Description: "下载文件", Usage: "wget [options] <url>", Example: "wget https://example.com/file.txt", Implemented: "system"},
		{Name: "ping", Category: "network", Description: "网络连通性测试", Usage: "ping [options] <host>", Example: "ping example.com", Implemented: "system"},
		{Name: "sed", Category: "text", Description: "文本替换", Usage: "sed [options] <expression> <file>", Example: "sed 's/old/new/g' file.txt", Implemented: "system"},
		{Name: "awk", Category: "text", Description: "文本处理", Usage: "awk [options] <pattern> <file>", Example: "awk '{print $1}' file.txt", Implemented: "system"},
		{Name: "cut", Category: "text", Description: "文本截取", Usage: "cut [options] <file>", Example: "cut -d',' -f1 file.csv", Implemented: "system"},
		{Name: "sort", Category: "text", Description: "文本排序", Usage: "sort [options] <file>", Example: "sort file.txt", Implemented: "system"},
		{Name: "wc", Category: "text", Description: "文本统计", Usage: "wc [options] <file>", Example: "wc -l file.txt", Implemented: "system"},
		{Name: "find", Category: "file", Description: "查找文件", Usage: "find [options] <path> <expression>", Example: "find . -name '*.go'", Implemented: "system"},
		{Name: "du", Category: "file", Description: "磁盘使用统计", Usage: "du [options] <path>", Example: "du -h .", Implemented: "system"},
		{Name: "top", Category: "system", Description: "实时进程监控", Usage: "top [options]", Example: "top", Implemented: "system"},
		{Name: "md5", Category: "encoding", Description: "MD5 哈希", Usage: "md5 [options] <file>", Example: "md5 file.txt", Implemented: "system"},
		{Name: "sha256", Category: "encoding", Description: "SHA256 哈希", Usage: "sha256 [options] <file>", Example: "sha256 file.txt", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}
