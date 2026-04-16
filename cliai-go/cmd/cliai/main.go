package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/pkgcmd"
)

var version = "0.1.0"

func main() {
	// 检查是否直接调用内置命令（不通过 cobra）
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		args := os.Args[2:]

		switch cmd {
		case "base64":
			if err := pkgcmd.ExecBase64(args); err != nil {
				fmt.Fprintln(os.Stderr, "错误:", err)
				os.Exit(1)
			}
			return
		case "curl":
			if err := pkgcmd.ExecCurl(args); err != nil {
				fmt.Fprintln(os.Stderr, "错误:", err)
				os.Exit(1)
			}
			return
		case "jq":
			if err := pkgcmd.ExecJq(args); err != nil {
				fmt.Fprintln(os.Stderr, "错误:", err)
				os.Exit(1)
			}
			return
		}
	}

	// 其他命令走 cobra
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
