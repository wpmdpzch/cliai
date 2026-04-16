package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Window TUI 窗口
type Window struct {
	width  int
	height int
	mode   string
	output []string
}

// NewWindow 创建 TUI 窗口
func NewWindow() *Window {
	return &Window{
		width:  100,
		height: 30,
		mode:   "CLI",
		output: []string{},
	}
}

// Run 运行窗口
func (w *Window) Run() error {
	// 设置终端
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	defer func() {
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
	}()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\033[2J\033[H")
	w.draw()

	for {
		// 光标到输入位置
		fmt.Printf("\033[%d;1H", w.height-2)
		fmt.Print("                                                        ")
		fmt.Printf("\033[%d;1H", w.height-2)
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			break
		}

		if input == "mode" {
			w.cycleMode()
			fmt.Print("\033[2J\033[H")
			w.draw()
			continue
		}

		w.output = append(w.output, fmt.Sprintf("[%s]$ %s", w.mode, input))
		w.output = append(w.output, fmt.Sprintf("→ 执行: %s", input))
		w.output = append(w.output, "")

		fmt.Print("\033[2J\033[H")
		w.draw()
	}

	fmt.Println("再见!")
	return nil
}

func (w *Window) cycleMode() {
	switch w.mode {
	case "CLI":
		w.mode = "PLAN"
	case "PLAN":
		w.mode = "BUILD"
	case "BUILD":
		w.mode = "CLI"
	}
}

func (w *Window) draw() {
	border := strings.Repeat("─", w.width-2)
	title := fmt.Sprintf(" CLI-AI v0.1.0  [%s] ", w.mode)
	hint := " Tab: 切换模式 | exit: 退出 "

	fmt.Println("┌" + border + "┐")
	fmt.Printf("│%s%-"+fmt.Sprintf("%d", w.width-2-len(title))+"s│\n", title, "")
	fmt.Printf("│%"+fmt.Sprintf("%d", w.width-2-len(hint))+"s%s│\n", "", hint)
	fmt.Println("├" + border + "┤")

	areaHeight := w.height - 10
	start := 0
	if len(w.output) > areaHeight {
		start = len(w.output) - areaHeight
	}

	for i := start; i < len(w.output); i++ {
		line := w.output[i]
		if len(line) > w.width-4 {
			line = line[:w.width-7] + "..."
		}
		fmt.Printf("│ %-"+fmt.Sprintf("%d", w.width-4)+"s │\n", line)
	}

	for i := len(w.output); i < areaHeight; i++ {
		fmt.Printf("│%"+fmt.Sprintf("%d", w.width-2)+"s │\n", "")
	}

	fmt.Println("└" + border + "┘")
}

// NewTUICommand 返回 TUI 窗口命令
func NewTUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "启动 TUI 窗口模式",
		Long:  `启动类似 OpenCode 的 TUI 窗口界面，支持 Tab 切换 CLI/PLAN/BUILD 模式。`,
		Run: func(cmd *cobra.Command, args []string) {
			window := NewWindow()
			window.Run()
		},
	}
}
