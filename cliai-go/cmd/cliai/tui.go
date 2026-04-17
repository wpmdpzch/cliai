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

// 配色方案 - 深色主题
const (
	// 主色调 - 蓝灰色系
	colorBg         = tcell.Color235   // 深灰背景 #1d1d1d
	colorBgLight    = tcell.Color238   // 浅灰 #2d2d2d
	colorBgInput    = tcell.Color234   // 输入框背景 #171717
	colorBorder     = tcell.Color240   // 边框灰 #3d3d3d
	colorText       = tcell.Color255   // 白色文字
	colorTextDim    = tcell.Color245   // 暗淡文字 #8a8a8a
	
	// 模式颜色
	colorModeCLI    = tcell.Color39    // 亮青 #00afff
	colorModePLAN   = tcell.Color226   // 金黄 #ffff00  
	colorModeBUILD  = tcell.Color82    // 亮绿 #00ff00
	
	// 文字颜色
	colorCmd       = tcell.Color252    // 浅灰命令文字 #c0c0c0
	colorError     = tcell.Color203    // 粉红错误 #ff6b6b
	colorSuccess   = tcell.Color114    // 绿色成功 #8aff6b
	colorInfo      = tcell.Color75     // 青色信息 #00afff
)

func NewTUI() *TUI {
	return &TUI{
		mode: core.ModeCLI,
	}
}

// TUI 主应用
type TUI struct {
	app        *tview.Application
	mode       core.Mode
	input      *tview.InputField
	output     *tview.TextView
	modeText   *tview.TextView
	header     *tview.Flex
	headerBg   *tview.TextView  // 标题背景色块
	headerVer  *tview.TextView  // 版本文字
}

func (t *TUI) Run() error {
	t.app = tview.NewApplication()

	// 设置全局样式
	tview.Styles.PrimitiveBackgroundColor = colorBg
	tview.Styles.PrimaryTextColor = colorText

	// 模式指示器 - 简洁条状
	t.modeText = tview.NewTextView()
	t.modeText.SetDynamicColors(true)
	t.modeText.SetBackgroundColor(colorBgLight)
	t.modeText.SetText(t.formatModeBar())

	// 输出区域 - 深色背景，浅色文字
	t.output = tview.NewTextView()
	t.output.SetDynamicColors(true)
	t.output.SetBackgroundColor(colorBg)
	t.output.SetTextColor(colorCmd)
	t.output.SetScrollable(true)
	t.output.SetWordWrap(true)
	t.output.SetBorder(false) // 无边框，更简洁
	t.output.SetText("[dim]欢迎使用 CLI-AI v0.1.0[::-]\n\n" +
		"[dim]功能:[::-]\n" +
		"  • 输入命令直接执行\n" +
		"  • 输入自然语言，AI 解析\n" +
		"  • [yellow]Tab[::-] 切换模式\n" +
		"  • 输入 [yellow]exit[::-] 退出\n\n" +
		"[dim]────────────────────────────[::-]")

	// 输入框 - 底部单行样式
	t.input = tview.NewInputField()
	t.input.SetBackgroundColor(colorBgInput)
	t.input.SetPlaceholderTextColor(colorTextDim)
	t.input.SetLabel(t.formatInputLabel())
	t.input.SetPlaceholder("输入命令或自然语言...")
	t.input.SetBorder(true)
	t.input.SetBorderColor(colorBorder)
	t.input.SetBorderAttributes(tcell.AttrBold)
	t.input.SetFieldBackgroundColor(colorBgInput) // 输入框内部背景
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

	// 顶部标题栏 - 分解为独立组件便于更新
	t.headerBg = tview.NewTextView()
	t.headerBg.SetTextAlign(tview.AlignLeft).
		SetText(" CLI-AI ").
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(colorModeCLI)

	t.headerVer = tview.NewTextView()
	t.headerVer.SetText(" v0.1.0 ").
		SetTextColor(colorTextDim).
		SetBackgroundColor(colorModeCLI)

	t.header = tview.NewFlex().
		AddItem(t.headerBg, 10, 0, false).
		AddItem(t.headerVer, 8, 0, false).
		AddItem(t.modeText, 0, 1, false)

	// 主布局 - 三行
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 0, false).
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
	var color tcell.Color
	var modeName string
	switch t.mode {
	case core.ModeCLI:
		color = colorModeCLI
		modeName = "CLI"
	case core.ModePlan:
		color = colorModePLAN
		modeName = "PLAN"
	case core.ModeBuild:
		color = colorModeBUILD
		modeName = "BUILD"
	}
	r, g, b := color.RGB()
	hexStr := fmt.Sprintf("%02x%02x%02x", r, g, b)
	return fmt.Sprintf(" [::b]%s%s[::-]::[::-] ", hexStr, modeName)
}

func (t *TUI) formatInputLabel() string {
	var color tcell.Color
	var modeName string
	switch t.mode {
	case core.ModeCLI:
		color = colorModeCLI
		modeName = "CLI"
	case core.ModePlan:
		color = colorModePLAN
		modeName = "PLAN"
	case core.ModeBuild:
		color = colorModeBUILD
		modeName = "BUILD"
	}
	r, g, b := color.RGB()
	hexStr := fmt.Sprintf("%02x%02x%02x", r, g, b)
	return fmt.Sprintf("[%s]%s>[-] ", hexStr, modeName)
}

func (t *TUI) updateMode() {
	// 更新模式文字
	t.modeText.SetText(t.formatModeBar())
	t.input.SetLabel(t.formatInputLabel())

	// 更新模式栏背景色
	switch t.mode {
	case core.ModeCLI:
		t.headerBg.SetBackgroundColor(colorModeCLI)
		t.headerVer.SetBackgroundColor(colorModeCLI)
	case core.ModePlan:
		t.headerBg.SetBackgroundColor(colorModePLAN)
		t.headerVer.SetBackgroundColor(colorModePLAN)
	case core.ModeBuild:
		t.headerBg.SetBackgroundColor(colorModeBUILD)
		t.headerVer.SetBackgroundColor(colorModeBUILD)
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

	// 记录命令 - 带模式前缀
	t.appendOutput(fmt.Sprintf("\n[dim]>%s>[-] %s\n", t.mode, cmd))

	// 执行命令
	result := pkgcmd.ExecCommand(cmd)
	if result.Error != nil {
		t.appendOutput(fmt.Sprintf("[red]✗ 错误:[-] %v\n", result.Error))
	} else if result.Output != "" {
		t.appendOutput(fmt.Sprintf("%s", result.Output))
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
