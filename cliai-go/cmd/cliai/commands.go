package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/pkgcmd"
)

var helpText = `
 CLI-AI v0.1.0 - 让命令行会思考

 用法:
   cliai [命令] [参数]

 命令:
   repl      启动交互式 REPL（默认）
   tui       启动 TUI 窗口模式
   help      显示帮助信息
   version   显示版本
   commands  列出可用命令

 交互模式:
   Tab 键切换三种模式:
   - CLI:   直接执行命令
   - PLAN:  预执行，只读操作
   - BUILD: 直接操作，可写文件

 示例:
   cliai repl
   cliai tui
   cliai help
   cliai commands
   cliai curl https://example.com
   cliai jq '.name' data.json

 配置文件:
   ~/.cliai/config.yaml

 内置命令包 (v0.1):
   network: curl
   text:   jq, grep, cat, sed, awk, cut, sort, wc
   file:   ls, head, tail, find, du, diff
   system: ps, top, df, free, du
   encoding: base64, md5sum, sha256sum

 文档:
   README.md   - 项目说明
   ROADMAP.md - 开发路线图
`

func NewHelpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "显示帮助信息",
		Long:  helpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(helpText)
		},
	}
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CLI-AI v0.1.0")
		},
	}
}

func NewCommandsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "列出可用命令",
		Run: func(cmd *cobra.Command, args []string) {
			cmds := pkgcmd.List()

			// 按类别分组
			categories := map[string][]*pkgcmd.Command{
				"network":  {},
				"text":     {},
				"file":     {},
				"system":   {},
				"encoding": {},
				"shell":    {},
			}

			for _, c := range cmds {
				if cat, ok := categories[c.Category]; ok {
					cat = append(cat, c)
					categories[c.Category] = cat
				}
			}

			// 颜色定义
			const (
				cyan    = "\033[36m"
				green   = "\033[32m"
				yellow  = "\033[33m"
				dim     = "\033[2m"
				reset   = "\033[0m"
			)

			fmt.Printf("%s 可用命令 (共 %d 个)%s\n\n", cyan, len(cmds), reset)

			catNames := map[string]string{
				"network":  "🌐 网络",
				"text":     "📝 文本",
				"file":     "📁 文件",
				"system":   "💻 系统",
				"encoding": "🔐 编码",
				"shell":    "🐚 Shell",
			}

			for cat, cmds := range categories {
				if len(cmds) > 0 {
					fmt.Printf("%s%s%s\n", cyan, catNames[cat], reset)
					fmt.Println(strings.Repeat("─", 40))
				for _, c := range cmds {
					impl := dim + "[sys]" + reset
					if c.Implemented == "go" {
						impl = green + "[go] " + reset
					} else if c.Implemented == "builtin" {
						impl = yellow + "[builtin]" + reset
					}
					fmt.Printf("  %-12s %s %s\n", c.Name, c.Description, impl)
				}
					fmt.Println()
				}
			}

			fmt.Printf("%s提示:%s 使用 %s cliai help <命令>%s 查看命令详情\n", 
				dim, reset, cyan, reset)
			fmt.Printf("%s      %s cliai <命令> --help%s 查看子命令帮助\n",
				dim, cyan, reset)
		},
	}
}
