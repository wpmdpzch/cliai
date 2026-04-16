package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "cliai",
		Short: "CLI-AI: 让命令行会思考",
		Long: `CLI-AI - 自然语言 CLI 工具
用自然语言控制命令行，支持内置命令包，跨平台开箱即用。`,
		Version: version,
	}

	rootCmd.AddCommand(NewREPLCommand())
	rootCmd.AddCommand(NewHelpCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewCommandsCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
