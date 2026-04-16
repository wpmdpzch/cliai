package builtin

import (
	"encoding/json"
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
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Usage       string   `json:"usage"`
	Example     string   `json:"example"`
	Implemented string   `json:"implemented"` // "go" or "system"
	Dangerous   bool     `json:"dangerous"`
	Args        []string `json:"args,omitempty"`
}

// NewCommandSet 创建命令集
func NewCommandSet() *CommandSet {
	cs := &CommandSet{
		commands: make(map[string]*Command),
	}
	cs.loadCommands()
	return cs
}

// loadCommands 加载内置命令
func (cs *CommandSet) loadCommands() {
	// 尝试从 builtin/commands.json 加载
	if data, err := os.ReadFile("builtin/commands.json"); err == nil {
		var commands []Command
		if err := json.Unmarshal(data, &commands); err == nil {
			for i := range commands {
				cs.commands[commands[i].Name] = &commands[i]
			}
		}
	}

	// 注册 P0 命令
	cs.registerP0Commands()
}

// registerP0Commands 注册 P0 命令
func (cs *CommandSet) registerP0Commands() {
	commands := []Command{
		// Network
		{Name: "curl", Category: "network", Description: "发送 HTTP/HTTPS 请求", Usage: "curl [options] <url>", Example: "curl https://example.com", Implemented: "go"},

		// Text
		{Name: "jq", Category: "text", Description: "JSON 处理", Usage: "jq [options] <filter> <file>", Example: "jq '.data' file.json", Implemented: "go"},
		{Name: "grep", Category: "text", Description: "文本搜索", Usage: "grep [options] <pattern> <file>", Example: "grep 'error' log.txt", Implemented: "system"},
		{Name: "cat", Category: "text", Description: "查看文件内容", Usage: "cat <file>", Example: "cat readme.md", Implemented: "system"},

		// File
		{Name: "ls", Category: "file", Description: "列出目录", Usage: "ls [options] <dir>", Example: "ls -la", Implemented: "system"},
		{Name: "head", Category: "file", Description: "查看文件头部", Usage: "head [options] <file>", Example: "head -n 10 file.txt", Implemented: "system"},
		{Name: "tail", Category: "file", Description: "查看文件尾部", Usage: "tail [options] <file>", Example: "tail -n 10 file.txt", Implemented: "system"},

		// System
		{Name: "ps", Category: "system", Description: "查看进程", Usage: "ps [options]", Example: "ps aux", Implemented: "system"},
		{Name: "df", Category: "system", Description: "查看磁盘使用", Usage: "df [options]", Example: "df -h", Implemented: "system"},
		{Name: "free", Category: "system", Description: "查看内存使用", Usage: "free [options]", Example: "free -h", Implemented: "system"},

		// Encoding
		{Name: "base64", Category: "encoding", Description: "Base64 编解码", Usage: "base64 [options] <file>", Example: "base64 file.txt", Implemented: "go"},
	}

	for i := range commands {
		cs.commands[commands[i].Name] = &commands[i]
	}
}

// Exists 检查命令是否存在
func (cs *CommandSet) Exists(name string) bool {
	_, ok := cs.commands[name]
	return ok
}

// Get 获取命令
func (cs *CommandSet) Get(name string) *Command {
	return cs.commands[name]
}

// List 返回所有命令
func (cs *CommandSet) List() []*Command {
	result := make([]*Command, 0, len(cs.commands))
	for _, cmd := range cs.commands {
		result = append(result, cmd)
	}
	return result
}

// Exec 执行命令
func (cs *CommandSet) Exec(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	name := parts[0]
	cmd, ok := cs.commands[name]
	if !ok {
		// 尝试系统命令
		return cs.execSystem(input)
	}

	if cmd.Implemented == "go" {
		return cs.execGo(cmd, parts[1:])
	}

	return cs.execSystem(input)
}

// execGo 执行 Go 原生命令
func (cs *CommandSet) execGo(cmd *Command, args []string) error {
	// TODO: 实现 Go 原生命令
	return nil
}

// execSystem 执行系统命令
func (cs *CommandSet) execSystem(input string) error {
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
