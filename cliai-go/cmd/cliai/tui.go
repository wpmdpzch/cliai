package main

import (
	"fmt"
	"os"
	"strings"

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
	colorHint      = tcell.Color238    // 提示文字背景
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
	headerBg   *tview.TextView
	headerVer  *tview.TextView
	suggestion *tview.TextView  // 提示区域
	history    []string         // 命令历史
	historyIdx int              // 历史索引
}

func (t *TUI) Run() error {
	t.app = tview.NewApplication()

	// 设置全局样式
	tview.Styles.PrimitiveBackgroundColor = colorBg
	tview.Styles.PrimaryTextColor = colorText

	// 初始化历史
	t.history = []string{}
	t.historyIdx = -1

	// 模式指示器
	t.modeText = tview.NewTextView()
	t.modeText.SetDynamicColors(true)
	t.modeText.SetBackgroundColor(colorBgLight)
	t.modeText.SetText(t.formatModeBar())

	// 输出区域
	t.output = tview.NewTextView()
	t.output.SetDynamicColors(true)
	t.output.SetBackgroundColor(colorBg)
	t.output.SetTextColor(colorCmd)
	t.output.SetScrollable(true)
	t.output.SetWordWrap(true)
	t.output.SetBorder(false)
	t.output.SetText("[dim]欢迎使用 CLI-AI v0.1.0[::-]\n\n" +
		"[dim]功能:[::-]\n" +
		"  • 输入命令直接执行\n" +
		"  • 输入自然语言，AI 解析\n" +
		"  • [yellow]Tab[::-] 切换模式\n" +
		"  • [cyan]Ctrl+I[::-] 补全命令\n" +
		"  • [cyan]↑↓[::-] 历史记录\n" +
		"  • 输入 [yellow]exit[::-] 退出\n\n" +
		"[dim]────────────────────────────[::-]")

	// 提示区域
	t.suggestion = tview.NewTextView()
	t.suggestion.SetDynamicColors(true)
	t.suggestion.SetBackgroundColor(colorHint)
	t.suggestion.SetTextColor(colorTextDim)
	t.suggestion.SetText("[dim]Ctrl+I 补全 | ↑↓ 历史 | Tab 切换模式[::-]")

	// 输入框
	t.input = tview.NewInputField()
	t.input.SetBackgroundColor(colorBgInput)
	t.input.SetPlaceholderTextColor(colorTextDim)
	t.input.SetLabel(t.formatInputLabel())
	t.input.SetPlaceholder("输入命令或自然语言...")
	t.input.SetBorder(true)
	t.input.SetBorderColor(colorBorder)
	t.input.SetBorderAttributes(tcell.AttrBold)
	t.input.SetFieldBackgroundColor(colorBgInput)
	t.input.SetAutocompleteFunc(t.getAutocomplete)

	// 输入框事件
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

	// 主布局 - 四行
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 0, false).
		AddItem(t.output, 0, 1, false).
		AddItem(t.suggestion, 1, 0, false).
		AddItem(t.input, 3, 0, false)

	t.app.SetRoot(flex, true)
	t.app.SetFocus(t.input)

	// 全局 Ctrl+C 退出
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			t.app.Stop()
			return nil
		}
		return event
	})

	return t.app.Run()
}

// getAutocomplete 获取自动补全候选项
func (t *TUI) getAutocomplete(currentText string) []string {
	// 解析当前输入
	parts := strings.Fields(currentText)
	if len(parts) == 0 {
		return nil
	}

	// 如果只有命令名，匹配命令
	if len(parts) == 1 {
		prefix := strings.ToLower(parts[0])
		var matches []string
		for _, cmd := range pkgcmd.List() {
			if strings.HasPrefix(strings.ToLower(cmd.Name), prefix) {
				matches = append(matches, cmd.Name+" - "+cmd.Description)
			}
		}
		return matches
	}

	return nil
}

// updateSuggestion 更新命令提示
func (t *TUI) updateSuggestion() {
	if t.app == nil {
		return
	}
	text := t.input.GetText()
	parts := strings.Fields(text)

	if len(parts) == 0 {
		t.app.QueueUpdate(func() {
			t.suggestion.SetText("[dim]Ctrl+I 补全 | ↑↓ 历史 | Tab 切换模式 | Ctrl+H 帮助[::-]")
		})
		return
	}

	// 查找匹配的命令
	cmdName := parts[0]
	if cmd := pkgcmd.Get(cmdName); cmd != nil {
		t.app.QueueUpdate(func() {
			t.suggestion.SetText(fmt.Sprintf("[cyan]%s[::-] %s [dim]例: %s[::-]",
				cmd.Name, cmd.Description, cmd.Example))
		})
	} else {
		// 查找以该字符串开头的命令
		var matches []string
		for _, cmd := range pkgcmd.List() {
			if strings.HasPrefix(strings.ToLower(cmd.Name), strings.ToLower(cmdName)) {
				matches = append(matches, cmd.Name)
			}
		}
		if len(matches) > 0 {
			t.app.QueueUpdate(func() {
				t.suggestion.SetText(fmt.Sprintf("[dim]候选:[::-] %s", strings.Join(matches, ", ")))
			})
		} else {
			t.app.QueueUpdate(func() {
				t.suggestion.SetText("[dim]未找到匹配命令 | Ctrl+I 补全[::-]")
			})
		}
	}
}

// navigateHistory 浏览历史
func (t *TUI) navigateHistory(dir int) {
	if len(t.history) == 0 {
		return
	}

	newIdx := t.historyIdx + dir

	if newIdx < 0 {
		// 到达最旧，恢复到当前输入
		newIdx = -1
		t.input.SetText("")
		return
	}

	if newIdx >= len(t.history) {
		// 到达最新，恢复空
		newIdx = len(t.history) - 1
	}

	t.historyIdx = newIdx
	t.input.SetText(t.history[len(t.history)-1-newIdx])
}

// showCommandHelp 显示命令帮助
func (t *TUI) showCommandHelp() {
	text := t.input.GetText()
	parts := strings.Fields(text)

	if len(parts) == 0 {
		t.appendOutput("[yellow]帮助:[::-] 输入命令名后按 Ctrl+H 查看该命令的帮助\n")
		return
	}

	cmdName := parts[0]
	if cmd := pkgcmd.Get(cmdName); cmd != nil {
		helpText := fmt.Sprintf("\n[yellow]命令:[::-] %s\n"+
			"[cyan]用法:[::-] %s\n"+
			"[green]说明:[::-] %s\n"+
			"[dim]示例:[::-] %s\n"+
			"[dim]类型:[::-] %s\n",
			cmd.Name, cmd.Usage, cmd.Description, cmd.Example, cmd.Implemented)
		t.appendOutput(helpText)
	} else {
		t.appendOutput(fmt.Sprintf("[red]未知命令:[::-] %s\n", cmdName))
	}
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
	t.modeText.SetText(t.formatModeBar())
	t.input.SetLabel(t.formatInputLabel())

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

	// 添加到历史
	t.history = append(t.history, cmd)
	t.historyIdx = -1 // 重置历史索引

	// 记录命令
	t.appendOutput(fmt.Sprintf("\n[dim]>%s>[::-] %s\n", t.mode, cmd))

	// 处理特殊命令
	if cmd == "help" {
		t.showFullHelp()
		return
	}

	// 执行命令
	result := pkgcmd.ExecCommand(cmd)
	if result.Error != nil {
		t.appendOutput(fmt.Sprintf("[red]✗ 错误:[-] %v\n", result.Error))
	} else if result.Output != "" {
		t.appendOutput(fmt.Sprintf("%s", result.Output))
	}
}

// showFullHelp 显示完整帮助
func (t *TUI) showFullHelp() {
	help := `
[cyan]═══════════════════════════════════════════[::-]
[cyan]  CLI-AI v0.1.0 - 让命令行会思考[::-]
[cyan]═══════════════════════════════════════════[::-]

[yellow]用法:[::-]
  cliai [命令] [参数]

[yellow]快捷键:[::-]
  [cyan]Tab[::-]      切换模式 (CLI/PLAN/BUILD)
  [cyan]Ctrl+I[::-]   命令补全
  [cyan]↑/↓[::-]      历史记录
  [cyan]Ctrl+H[::-]   命令帮助
  [cyan]Ctrl+C[::-]   退出

[yellow]模式:[::-]
  [blue]CLI[::-]   直接执行命令
  [yellow]PLAN[::-]  预执行模式（只读）
  [green]BUILD[::-] 直接操作（可写）

[yellow]内置命令 (73个):[::-]
  🌐 网络:    curl, wget, ping, ip, ss, netstat, nslookup, dig...
  📝 文本:    jq, grep, cat, sed, awk, cut, sort, wc, diff...
  📁 文件:    ls, find, head, tail, ln, chmod, chown, tar...
  💻 系统:    ps, top, df, free, kill, htop, systemctl...
  🔐 编码:    base64, md5sum, sha256sum
  🐚 Shell:   cd, pwd, echo, mkdir, touch, rm, cp, mv...

[yellow]示例:[::-]
  cliai curl https://example.com
  cliai jq '.name' data.json
  cliai ls -la /tmp
  cliai ps aux | grep nginx

[dim]提示:[::-] 输入自然语言，AI 会帮你解析成命令！

[cyan]═══════════════════════════════════════════[::-]
`
	t.appendOutput(help)
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
