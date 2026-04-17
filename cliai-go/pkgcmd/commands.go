package pkgcmd

import (
	"bytes"
	"fmt"
	"os"
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
	registerP2Commands()
	registerP3Commands()
	registerP4Commands()
	registerP5Commands()
	registerShellBuiltins()
}

// registerShellBuiltins 注册 shell 内建命令
func registerShellBuiltins() {
	cmds := []Command{
		{Name: "cd", Category: "shell", Description: "切换目录", Usage: "cd <目录>", Example: "cd /tmp", Implemented: "builtin"},
		{Name: "echo", Category: "shell", Description: "输出文本", Usage: "echo [选项] <文本>", Example: "echo hello world", Implemented: "builtin"},
		{Name: "pwd", Category: "shell", Description: "显示当前目录", Usage: "pwd", Example: "pwd", Implemented: "builtin"},
		{Name: "mkdir", Category: "shell", Description: "创建目录", Usage: "mkdir [选项] <目录>", Example: "mkdir -p path/to/dir", Implemented: "builtin"},
		{Name: "touch", Category: "shell", Description: "创建文件", Usage: "touch <文件>", Example: "touch readme.md", Implemented: "builtin"},
		{Name: "rm", Category: "shell", Description: "删除文件", Usage: "rm [选项] <文件>", Example: "rm file.txt", Dangerous: true, Implemented: "builtin"},
		{Name: "cp", Category: "shell", Description: "复制文件", Usage: "cp <源> <目标>", Example: "cp a.txt b.txt", Implemented: "builtin"},
		{Name: "mv", Category: "shell", Description: "移动/重命名文件", Usage: "mv <源> <目标>", Example: "mv a.txt b.txt", Implemented: "builtin"},
		{Name: "clear", Category: "shell", Description: "清屏", Usage: "clear", Example: "clear", Implemented: "builtin"},
		{Name: "exit", Category: "shell", Description: "退出", Usage: "exit [代码]", Example: "exit 0", Implemented: "builtin"},
		{Name: "which", Category: "shell", Description: "查找命令位置", Usage: "which <命令>", Example: "which ls", Implemented: "builtin"},
		{Name: "history", Category: "shell", Description: "查看命令历史", Usage: "history", Example: "history", Implemented: "builtin"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
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

	// Go 原生实现或内置命令都走 execGoCmd
	if cmd.Implemented == "go" || cmd.Implemented == "builtin" {
		return execGoCmd(cmd, parts[1:])
	}

	return execSystemCmd(input)
}

// execGoCmd 执行 Go 原生命令
func execGoCmd(cmd *Command, args []string) *ExecResult {
	// 先尝试 Go 原生实现
	switch cmd.Name {
	case "base64":
		var buf bytes.Buffer
		err := ExecBase64ToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	case "curl":
		var buf bytes.Buffer
		err := ExecCurlToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	case "jq":
		var buf bytes.Buffer
		err := ExecJqToWriter(args, &buf)
		return &ExecResult{Output: buf.String(), Error: err}
	}

	// 如果不是 Go 原生命令但标记为 builtin，调用 builtin 处理
	if cmd.Implemented == "builtin" {
		return ExecGoCmdInternal(cmd.Name, args, nil)
	}

	return &ExecResult{Output: "", Error: fmt.Errorf("未实现的命令: %s", cmd.Name)}
}

// execSystemCmd 执行系统命令
func execSystemCmd(input string) *ExecResult {
	args := strings.Fields(input)
	if len(args) == 0 {
		return &ExecResult{Output: "", Error: nil}
	}

	cmd := exec.Command(args[0], args[1:]...)
	
	// 检查是否有 stdin 输入
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// stdin 有数据，传递给它
		cmd.Stdin = os.Stdin
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// 如果有 stderr 输出但没有错误，仍然返回 stderr 内容（警告）
	if err != nil && stderr.Len() > 0 {
		return &ExecResult{Output: "", Error: fmt.Errorf("%s", strings.TrimSpace(stderr.String()))}
	}
	
	return &ExecResult{Output: stdout.String(), Error: err}
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
		{Name: "md5sum", Category: "encoding", Description: "MD5 哈希", Usage: "md5sum [options] <file>", Example: "md5sum file.txt", Implemented: "system"},
		{Name: "sha256sum", Category: "encoding", Description: "SHA256 哈希", Usage: "sha256sum [options] <file>", Example: "sha256sum file.txt", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}

// registerP2Commands 注册 P2 网络工具命令
func registerP2Commands() {
	cmds := []Command{
		// 网络工具 (net-tools / iproute2)
		{Name: "ip", Category: "network", Description: "IP 地址和路由管理", Usage: "ip [选项] <命令>", Example: "ip addr show", Implemented: "system"},
		{Name: "ss", Category: "network", Description: "Socket 统计", Usage: "ss [选项] <过滤器>", Example: "ss -tuln", Implemented: "system"},
		{Name: "netstat", Category: "network", Description: "网络状态查看", Usage: "netstat [选项]", Example: "netstat -tuln", Implemented: "system"},
		{Name: "ifconfig", Category: "network", Description: "网络接口配置", Usage: "ifconfig [接口] [选项]", Example: "ifconfig eth0", Implemented: "system"},
		{Name: "route", Category: "network", Description: "路由表管理", Usage: "route [选项]", Example: "route -n", Implemented: "system"},
		{Name: "arp", Category: "network", Description: "ARP 缓存管理", Usage: "arp [选项]", Example: "arp -a", Implemented: "system"},
		{Name: "nslookup", Category: "network", Description: "DNS 查询", Usage: "nslookup <域名>", Example: "nslookup example.com", Implemented: "system"},
		{Name: "dig", Category: "network", Description: "DNS 详细查询", Usage: "dig [选项] <域名>", Example: "dig example.com", Implemented: "system"},
		{Name: "host", Category: "network", Description: "DNS lookup", Usage: "host <域名>", Example: "host example.com", Implemented: "system"},
		{Name: "traceroute", Category: "network", Description: "路由追踪", Usage: "traceroute [选项] <主机>", Example: "traceroute example.com", Implemented: "system"},
		{Name: "tracepath", Category: "network", Description: "路径追踪", Usage: "tracepath <主机>", Example: "tracepath example.com", Implemented: "system"},
		{Name: "nc", Category: "network", Description: "网络瑞士军刀", Usage: "nc [选项] <主机> <端口>", Example: "nc -l 8080", Implemented: "system"},
		{Name: "telnet", Category: "network", Description: "远程登录测试", Usage: "telnet <主机> <端口>", Example: "telnet example.com 80", Implemented: "system"},
		{Name: "nmap", Category: "network", Description: "端口扫描", Usage: "nmap [选项] <目标>", Example: "nmap -sV localhost", Implemented: "system"},
		{Name: "tcpdump", Category: "network", Description: "抓包分析", Usage: "tcpdump [选项] [表达式>", Example: "tcpdump -i eth0", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}

// registerP3Commands 注册 P3 文件操作命令
func registerP3Commands() {
	cmds := []Command{
		{Name: "ln", Category: "file", Description: "创建链接", Usage: "ln [选项] <源> <目标>", Example: "ln -s file.txt link", Implemented: "system"},
		{Name: "chmod", Category: "file", Description: "权限管理", Usage: "chmod [选项] <模式> <文件>", Example: "chmod 755 script.sh", Implemented: "system"},
		{Name: "chown", Category: "file", Description: "所有者管理", Usage: "chown [选项] <用户>:<组> <文件>", Example: "chown user:group file", Implemented: "system"},
		{Name: "tar", Category: "file", Description: "打包/解压", Usage: "tar [选项] <文件>", Example: "tar -czvf file.tar.gz dir/", Implemented: "system"},
		{Name: "zip", Category: "file", Description: "压缩", Usage: "zip [选项] <压缩包> <文件>", Example: "zip -r file.zip dir/", Implemented: "system"},
		{Name: "unzip", Category: "file", Description: "解压", Usage: "unzip [选项] <压缩包>", Example: "unzip file.zip", Implemented: "system"},
		{Name: "rsync", Category: "file", Description: "同步", Usage: "rsync [选项] <源> <目标>", Example: "rsync -avz src/ dst/", Implemented: "system"},
		{Name: "tree", Category: "file", Description: "目录树展示", Usage: "tree [选项] <目录>", Example: "tree -L 2", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}

// registerP4Commands 注册 P4 文本处理命令
func registerP4Commands() {
	cmds := []Command{
		{Name: "less", Category: "text", Description: "分页查看", Usage: "less <文件>", Example: "less file.txt", Implemented: "system"},
		{Name: "uniq", Category: "text", Description: "去重", Usage: "uniq [选项] <文件>", Example: "uniq file.txt", Implemented: "system"},
		{Name: "diff", Category: "text", Description: "对比文件", Usage: "diff [选项] <文件1> <文件2>", Example: "diff a.txt b.txt", Implemented: "system"},
		{Name: "paste", Category: "text", Description: "合并文件", Usage: "paste [选项] <文件列表>", Example: "paste -d',' f1 f2", Implemented: "system"},
		{Name: "tr", Category: "text", Description: "字符转换", Usage: "tr [选项] <字符集1> <字符集2>", Example: "tr 'a-z' 'A-Z'", Implemented: "system"},
		{Name: "xargs", Category: "text", Description: "参数构建", Usage: "xargs [选项] <命令>", Example: "echo 'a b c' | xargs -n1", Implemented: "system"},
		{Name: "tee", Category: "text", Description: "管道复制", Usage: "tee [选项] <文件>", Example: "echo 'test' | tee file.txt", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}

// registerP5Commands 注册 P5 系统管理命令
func registerP5Commands() {
	cmds := []Command{
		{Name: "kill", Category: "system", Description: "终止进程", Usage: "kill [选项] <PID>", Example: "kill -9 1234", Implemented: "system"},
		{Name: "pkill", Category: "system", Description: "按名杀进程", Usage: "pkill <进程名>", Example: "pkill firefox", Implemented: "system"},
		{Name: "htop", Category: "system", Description: "增强进程监控", Usage: "htop [选项]", Example: "htop", Implemented: "system"},
		{Name: "systemctl", Category: "system", Description: "服务管理", Usage: "systemctl [命令] <服务>", Example: "systemctl status nginx", Implemented: "system"},
		{Name: "mount", Category: "system", Description: "挂载文件系统", Usage: "mount [选项] <设备> <目录>", Example: "mount /dev/sdb1 /mnt", Implemented: "system"},
		{Name: "umount", Category: "system", Description: "卸载文件系统", Usage: "umount <目录>", Example: "umount /mnt", Implemented: "system"},
		{Name: "lsof", Category: "system", Description: "查看打开文件", Usage: "lsof [选项]", Example: "lsof -i :8080", Implemented: "system"},
		{Name: "strace", Category: "system", Description: "系统调用追踪", Usage: "strace [选项] <命令>", Example: "strace -f ls", Implemented: "system"},
	}
	for i := range cmds {
		commands[cmds[i].Name] = &cmds[i]
	}
}
