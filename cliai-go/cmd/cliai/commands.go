package main

import (
	"fmt"

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
  cliai commands
  cliai version

配置文件:
  ~/.cliai/config.yaml

内置命令包 (v0.1):
  network: curl
  text:   jq, grep, cat
  file:   ls, head, tail
  system: ps, df, free
  encoding: base64
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

			fmt.Println("可用命令:")
			fmt.Println()

			// 按类别分组
			categories := map[string][]*pkgcmd.Command{
				"network":  {},
				"text":     {},
				"file":     {},
				"system":   {},
				"encoding": {},
			}

			for _, c := range cmds {
				if cat, ok := categories[c.Category]; ok {
					cat = append(cat, c)
					categories[c.Category] = cat
				}
			}

			for cat, cmds := range categories {
				if len(cmds) > 0 {
					fmt.Printf("%s:\n", cat)
					for _, c := range cmds {
						impl := "sys"
						if c.Implemented == "go" {
							impl = "go"
						}
						fmt.Printf("  %-10s %s [%s]\n", c.Name, c.Description, impl)
					}
					fmt.Println()
				}
			}
		},
	}
}
