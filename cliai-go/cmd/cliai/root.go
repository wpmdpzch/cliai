package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// NewREPLCommand 返回 REPL 命令
func NewREPLCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repl",
		Short: "启动交互式 REPL",
		Long:  `启动 CLI-AI 交互式命令行界面，支持 Tab 切换 CLI/PLAN/BUILD 模式。`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := loadConfig(cfgFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
				os.Exit(1)
			}

			repl := NewREPL(cfg)
			if err := repl.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "REPL 错误: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// loadConfig 加载配置
func loadConfig(path string) (*Config, error) {
	if path == "" {
		path = os.ExpandEnv("$HOME/.cliai/config.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// 返回默认配置
		return DefaultConfig(), nil
	}

	return ParseConfig(data)
}
