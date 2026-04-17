package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/pkgcmd"
)

var version = "0.1.0"

//的颜色定义
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorDim    = "\033[2m"
	colorReset  = "\033[0m"
)

// 跨平台清屏
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[2J\033[H")
	}
}

// 打印欢迎信息
func printWelcome() {
	fmt.Printf("%s╭─────────────────────────────────────────%s\n", colorCyan, colorReset)
	fmt.Printf("%s│%s  CLI-AI %s v%s - 让命令行会思考%s\n", colorCyan, colorReset, colorYellow, version, colorReset)
	fmt.Printf("%s╰─────────────────────────────────────────%s\n", colorCyan, colorReset)
	fmt.Println()
}

// 判断命令是否可能涉及二进制文件
func isBinaryCommand(cmdName string) bool {
	binCmds := []string{"cat", "head", "tail", "less", "more"}
	for _, c := range binCmds {
		if cmdName == c {
			return true
		}
	}
	return false
}

// 判断输出是否可能是二进制内容
func isLikelyBinary(data string) bool {
	// 检查是否包含大量不可打印字符
	nonPrintable := 0
	for _, c := range data {
		if c < 32 && c != '\n' && c != '\t' {
			nonPrintable++
		}
	}
	// 如果超过 10% 是不可打印字符，认为是二进制
	return float64(nonPrintable)/float64(len(data)) > 0.1
}

func main() {
	// 检查是否直接调用内置命令（不通过 cobra）
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		args := os.Args[2:]

		// 帮助类参数
		if cmd == "-h" || cmd == "--help" {
			printWelcome()
			fmt.Println("用法: cliai [命令] [参数]")
			fmt.Println()
			fmt.Println("输入 cliai help 查看完整帮助")
			return
		}

		// 检查是否是注册的命令
		if pkgcmd.Exists(cmd) {
			// 构建完整命令字符串
			fullCmd := cmd
			for _, arg := range args {
				// 跳过颜色相关的参数
				if !strings.HasPrefix(arg, "--color") && !strings.HasPrefix(arg, "-c") {
					fullCmd += " " + arg
				}
			}
			// 执行命令
			result := pkgcmd.ExecCommand(fullCmd)
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "%s错误:%s %v\n", colorRed, colorReset, result.Error)
				os.Exit(1)
			}
			
			// 如果是 cat/head/tail 等命令，检查输出是否是二进制
			if isBinaryCommand(cmd) && isLikelyBinary(result.Output) {
				fmt.Fprintf(os.Stderr, "%s警告:%s 输出可能是二进制内容，已跳过显示\n", colorYellow, colorReset)
				fmt.Fprintf(os.Stderr, "%s提示:%s 使用 %s xxd%s 或 %s hexdump%s 查看二进制\n", 
					colorDim, colorReset, colorCyan, colorReset, colorCyan, colorReset)
				return
			}
			
			fmt.Print(result.Output)
			return
		}
	}

	// 无参数或未知参数时，显示帮助
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "help") {
		NewHelpCommand().Run(nil, nil)
		return
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
	rootCmd.AddCommand(NewTUICommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s错误:%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}
}
