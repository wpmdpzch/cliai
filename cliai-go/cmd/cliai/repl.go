package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wpmdpzch/cliai/core"
	"github.com/wpmdpzch/cliai/pkgcmd"
)

// REPL 命令行交互
type REPL struct {
	mode     core.Mode
	scanner  *bufio.Scanner
	aiEngine *core.AIEngine
}

func NewREPL() *REPL {
	r := &REPL{
		mode:    core.ModeCLI,
		scanner: bufio.NewScanner(os.Stdin),
	}
	r.aiEngine = core.NewAIEngine(nil)
	return r
}

// Run 启动 REPL
func (r *REPL) Run() error {
	fmt.Println("CLI-AI v0.1.0 - 让命令行会思考")
	fmt.Println("Tab 键切换模式: CLI → PLAN → BUILD")
	fmt.Println("输入 'exit' 退出")
	fmt.Println()

	for {
		fmt.Printf("[%s]$ ", r.mode)
		if !r.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(r.scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("再见!")
			break
		}

		if input == "mode" {
			r.cycleMode()
			continue
		}

		if err := r.handleInput(input); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return nil
}

// cycleMode 循环切换模式
func (r *REPL) cycleMode() {
	r.mode.Next()
	fmt.Printf("切换到 %s 模式\n", r.mode)
}

// handleInput 处理输入
func (r *REPL) handleInput(input string) error {
	// 检查是否是内置命令
	if pkgcmd.Exists(input) {
		result := pkgcmd.ExecCommand(input)
		if result.Error != nil {
			return result.Error
		}
		fmt.Print(result.Output)
		return nil
	}

	// AI 解析自然语言
	return r.aiEngine.Process(input, r.mode)
}
