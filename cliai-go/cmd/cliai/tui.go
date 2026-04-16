package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/core"
	"github.com/wpmdpzch/cliai/pkgcmd"
	"github.com/gdamore/tcell/v2"
)

// TUI 主应用
type TUI struct {
	app      *tview.Application
	mode     core.Mode
	input    *tview.InputField
	output   *tview.TextView
	modeText *tview.TextView
}

func NewTUI() *TUI {
	return &TUI{
		mode: core.ModeCLI,
	}
}

func (t *TUI) Run() error {
	t.app = tview.NewApplication()

	// 模式指示器
	t.modeText = tview.NewTextView().
		SetDynamicColors(true).
		SetText(t.formatModeBar())

	// 输出区域
	t.output = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetText("CLI-AI v0.1.0 - 让命令行会思考\n\n输入命令或自然语言，按 Enter 执行\n按 Tab 切换模式\n")

	// 输入框
	t.input = tview.NewInputField().
		SetLabel(t.formatInputLabel()).
		SetPlaceholder("输入命令或自然语言...").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				t.executeCommand()
			}
		})

	// Tab 切换模式
	t.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			t.mode.Next()
			t.updateMode()
			return nil
		}
		return event
	})

	// 布局: 顶部模式条 + 中间输出 + 底部输入
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.modeText, 1, 0, false).
		AddItem(t.output, 0, 1, false).
		AddItem(t.input, 1, 0, false)

	t.app.SetRoot(flex, true)
	t.app.SetFocus(t.input)

	return t.app.Run()
}

func (t *TUI) formatModeBar() string {
	return fmt.Sprintf("[::b]%s[::-] 模式  |  Tab: 切换  |  exit: 退出", t.mode)
}

func (t *TUI) formatInputLabel() string {
	return fmt.Sprintf("[%s]$ ", t.mode)
}

func (t *TUI) updateMode() {
	t.modeText.SetText(t.formatModeBar())
	t.input.SetLabel(t.formatInputLabel())
}

func (t *TUI) executeCommand() {
	cmd := t.input.GetText()
	t.input.SetText("")

	if cmd == "exit" || cmd == "quit" {
		t.app.Stop()
		return
	}

	// 记录命令
	t.appendOutput(fmt.Sprintf("\n[%s]$ %s\n", t.mode, cmd))

	// 执行命令
	if err := pkgcmd.ExecCommand(cmd); err != nil {
		t.appendOutput(fmt.Sprintf("执行失败: %v\n", err))
	}
}

func (t *TUI) appendOutput(s string) {
	t.output.SetText(t.output.GetText(false) + s)
}

// NewTUICommand 返回 TUI 命令
func NewTUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "启动 TUI 窗口模式",
		Long:  `启动类似 OpenCode 的 TUI 窗口界面，支持 Tab 切换 CLI/PLAN/BUILD 模式。`,
		Run: func(cmd *cobra.Command, args []string) {
			tui := NewTUI()
			if err := tui.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "TUI 错误: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
