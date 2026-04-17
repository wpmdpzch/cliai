package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/config"
)

// loadConfig 加载配置
func loadConfig() *config.Config {
	usr, err := user.Current()
	if err != nil {
		return config.DefaultConfig()
	}
	cfgPath := filepath.Join(usr.HomeDir, ".cliai", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return config.DefaultConfig()
	}
	return cfg
}

// NewREPLCommand 返回 REPL 命令
func NewREPLCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repl",
		Short: "启动交互式 REPL",
		Long:  `启动 CLI-AI 交互式命令行界面，支持 Tab 切换 CLI/PLAN/BUILD 模式。`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := loadConfig()
			repl := NewREPL(cfg)
			if err := repl.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "REPL 错误: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
