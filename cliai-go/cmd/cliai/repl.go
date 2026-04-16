package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wpmdpzch/cliai/pkgcmd"
	"github.com/wpmdpzch/cliai/core"
	"github.com/wpmdpzch/cliai/config"
)

// Mode REPL 模式
type Mode int

const (
	ModeCLI  Mode = iota // 直接执行
	ModePlan             // 预执行（只读）
	ModeBuild            // 直接操作（可写）
)

func (m Mode) String() string {
	switch m {
	case ModeCLI:
		return "CLI"
	case ModePlan:
		return "PLAN"
	case ModeBuild:
		return "BUILD"
	default:
		return "CLI"
	}
}

// REPL 命令行交互
type REPL struct {
	mode       Mode
	scanner    *bufio.Scanner
	commands   *builtin.CommandSet
	aiEngine   *core.AIEngine
	cfg        *config.Config
}

func NewREPL(cfg *config.Config) *REPL {
	return &REPL{
		mode:     ModeCLI,
		scanner:  bufio.NewScanner(os.Stdin),
		commands: builtin.NewCommandSet(),
		cfg:      cfg,
	}
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
	switch r.mode {
	case ModeCLI:
		r.mode = ModePlan
	case ModePlan:
		r.mode = ModeBuild
	case ModeBuild:
		r.mode = ModeCLI
	}
	fmt.Printf("切换到 %s 模式\n", r.mode)
}

// handleInput 处理输入
func (r *REPL) handleInput(input string) error {
	// 检查是否是内置命令
	if r.commands.Exists(input) {
		return r.commands.Exec(input)
	}

	// AI 解析自然语言
	return r.aiEngine.Process(input, r.mode)
}
