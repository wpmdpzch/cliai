package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/wpmdpzch/cliai/core"
	"github.com/wpmdpzch/cliai/pkgcmd"
)

func NewTUI() *TUI {
	return &TUI{
		mode: core.ModeCLI,
	}
}

// TUI 主应用
type TUI struct {
	app      *tview.Application
	mode     core.Mode
	input    *tview.InputField
	output   *tview.TextView
	modeText *tview.TextView
}

func (t *TUI) Run() error {
	t.app = tview.NewApplication()

	// 模式指示器
	t.modeText = tview.NewTextView()
	t.modeText.SetDynamicColors(true)
	t.modeText.SetText(t.formatModeBar())
	t.modeText.SetBackgroundColor(tcell.ColorTeal)

	// 输出区域
	t.output = tview.NewTextView()
	t.output.SetDynamicColors(true)
	t.output.SetScrollable(true)
	t.output.SetWordWrap(true)
	t.output.SetBorder(true)
	t.output.SetBorderColor(tcell.ColorTeal)
	t.output.SetBorderAttributes(tcell.AttrBold)
	t.output.SetTitle(" 输出 ")
	t.output.SetText("[::b]欢迎使用 CLI-AI v0.1.0[::-]\n\n[::b]功能:[::-]\n  • 输入命令直接执行\n  • 输入自然语言，AI 解析\n  • [yellow]Tab[::-] 切换模式\n  • 输入 [yellow]exit[::-] 退出\n\n[dim]按 Tab 切换模式体验颜色变化[::-]")

	// 输入框
	t.input = tview.NewInputField()
	t.input.SetLabel(t.formatInputLabel())
	t.input.SetPlaceholder("输入命令或自然语言...")
	t.input.SetBorder(true)
	t.input.SetBorderColor(tcell.ColorYellow)
	t.input.SetBorderAttributes(tcell.AttrBold)
	t.input.SetDoneFunc(func(key tcell.Key) {
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

	// 顶部标题栏
	header := tview.NewFlex().
		AddItem(tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetText(" CLI-AI ").
			SetDynamicColors(true).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorPurple), 20, 0, false).
		AddItem(t.modeText, 0, 1, false)

	// 主布局
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(t.output, 0, 1, false).
		AddItem(t.input, 3, 0, false)

	t.app.SetRoot(flex, true)
	t.app.SetFocus(t.input)

	// Ctrl+C 退出
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			t.app.Stop()
			return nil
		}
		return event
	})

	return t.app.Run()
}

func (t *TUI) formatModeBar() string {
	return fmt.Sprintf(" [::b]%s[::-] 模式 ", t.mode)
}

func (t *TUI) formatInputLabel() string {
	var color string
	switch t.mode {
	case core.ModeCLI:
		color = "[green]"
	case core.ModePlan:
		color = "[yellow]"
	case core.ModeBuild:
		color = "[red]"
	}
	return fmt.Sprintf("%s%s > [::-]", color, t.mode)
}

func (t *TUI) updateMode() {
	// 更新模式文字颜色
	t.modeText.SetText(t.formatModeBar())
	t.input.SetLabel(t.formatInputLabel())

	// 更新模式栏背景色
	switch t.mode {
	case core.ModeCLI:
		t.modeText.SetBackgroundColor(tcell.ColorGreen)
	case core.ModePlan:
		t.modeText.SetBackgroundColor(tcell.ColorYellow)
	case core.ModeBuild:
		t.modeText.SetBackgroundColor(tcell.ColorRed)
	}
}

func (t *TUI) executeCommand() {
	cmd := t.input.GetText()
	t.input.SetText("")

	if cmd == "exit" || cmd == "quit" {
		t.app.Stop()
		return
	}

	if cmd == "" {
		return
	}

	// 记录命令
	t.appendOutput(fmt.Sprintf("\n[dim]%s > [::-]%s\n", t.mode, cmd))

	// 执行命令并捕获输出
	result := pkgcmd.ExecCommand(cmd)
	if result.Error != nil {
		t.appendOutput(fmt.Sprintf("[red]错误: [::-]%v\n", result.Error))
	} else if result.Output != "" {
		t.appendOutput(fmt.Sprintf("%s\n", result.Output))
	}
}

func (t *TUI) appendOutput(s string) {
	t.output.SetText(t.output.GetText(false) + s)
	t.output.ScrollToEnd()
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
