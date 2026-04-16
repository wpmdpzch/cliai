package pkgcmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandSet 内置命令集
type CommandSet struct {
	commands map[string]*Command
}

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

// ExecCommand 执行任意命令
func ExecCommand(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	name := parts[0]
	cmd, ok := commands[name]
	if !ok {
		return execSystem(input)
	}

	if cmd.Implemented == "go" {
		return execGo(cmd, parts[1:])
	}

	return execSystem(input)
}

// execGo 执行 Go 原生命令
func execGo(cmd *Command, args []string) error {
	switch cmd.Name {
	case "base64":
		return ExecBase64(args)
	case "curl":
		return ExecCurl(args)
	case "jq":
		return ExecJq(args)
	default:
		return fmt.Errorf("未实现的 Go 命令: %s", cmd.Name)
	}
}

// execSystem 执行系统命令
func execSystem(input string) error {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
