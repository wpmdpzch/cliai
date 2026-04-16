package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewREPLCommand 返回 REPL 命令
func NewREPLCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repl",
		Short: "启动交互式 REPL",
		Long:  `启动 CLI-AI 交互式命令行界面，支持 Tab 切换 CLI/PLAN/BUILD 模式。`,
		Run: func(cmd *cobra.Command, args []string) {
			repl := NewREPL()
			if err := repl.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "REPL 错误: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
