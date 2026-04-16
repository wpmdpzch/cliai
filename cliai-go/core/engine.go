package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wpmdpzch/cliai/pkgcmd"
	"github.com/wpmdpzch/cliai/config"
)

// Mode 运行模式
type Mode int

const (
	ModeCLI  Mode = iota
	ModePlan
	ModeBuild
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

// AIEngine AI 引擎
type AIEngine struct {
	cfg      *config.Config
	commands *pkgcmd.CommandSet
	client   *http.Client
}

// NewAIEngine 创建 AI 引擎
func NewAIEngine(cfg *config.Config, commands *pkgcmd.CommandSet) *AIEngine {
	return &AIEngine{
		cfg:      cfg,
		commands: commands,
		client:   &http.Client{},
	}
}

// Process 处理输入
func (e *AIEngine) Process(input string, mode Mode) error {
	// 构建 prompt
	prompt := e.buildPrompt(input, mode)

	// 调用 AI
	response, err := e.callAI(prompt)
	if err != nil {
		return fmt.Errorf("AI 调用失败: %v", err)
	}

	// 解析响应
	actions, err := e.parseResponse(response)
	if err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 执行操作
	return e.executeActions(actions, mode)
}

// buildPrompt 构建提示
func (e *AIEngine) buildPrompt(input string, mode Mode) string {
	commands := e.commands.List()

	var cmdList strings.Builder
	cmdList.WriteString("可用命令:\n")
	for _, cmd := range commands {
		cmdList.WriteString(fmt.Sprintf("- %s: %s (%s)\n", cmd.Name, cmd.Description, cmd.Implemented))
	}

	systemPrompt := "你是一个命令行助手。用户输入自然语言，你需要生成相应的命令。\n"
	systemPrompt += cmdList.String()
	systemPrompt += "\n根据用户输入，生成要执行的命令。返回 JSON 格式：{\"commands\": [\"cmd1\", \"cmd2\", ...]}\n"

	if mode == ModePlan {
		systemPrompt = strings.Replace(systemPrompt, "执行命令", "只生成命令，不执行", 1)
	}

	return systemPrompt + "\n用户: " + input
}

// callAI 调用 AI API
func (e *AIEngine) callAI(prompt string) (string, error) {
	url := e.cfg.AI.BaseURL + "/chat/completions"

	reqBody := map[string]interface{}{
		"model": e.cfg.AI.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": e.cfg.AI.Temp,
		"max_tokens":  e.cfg.AI.MaxTokens,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.cfg.AI.APIKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误: %d - %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("无效的 API 响应")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无效的响应格式")
	}

	msg, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无效的消息格式")
	}

	content, ok := msg["content"].(string)
	if !ok {
		return "", fmt.Errorf("无效的内容格式")
	}

	return content, nil
}

// parseResponse 解析 AI 响应
func (e *AIEngine) parseResponse(response string) ([]string, error) {
	// 尝试提取 JSON
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		jsonStr := response[start : end+1]
		var result map[string][]string
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			if cmds, ok := result["commands"]; ok {
				return cmds, nil
			}
		}
	}

	// 降级：返回原始响应作为单条命令
	return []string{response}, nil
}

// executeActions 执行操作
func (e *AIEngine) executeActions(actions []string, mode Mode) error {
	for _, action := range actions {
		action = strings.TrimSpace(action)

		// 移除 markdown 代码块
		action = strings.Trim(action, "` \n")

		if mode == ModePlan {
			fmt.Printf("[PLAN] 只读模式，不执行: %s\n", action)
			continue
		}

		// 检查是否是危险命令
		if e.isDangerous(action) && e.cfg.Exec.ConfirmDangerous {
			fmt.Printf("⚠️ 危险命令: %s\n", action)
			fmt.Print("确认执行? (y/N): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("已取消")
				continue
			}
		}

		// 执行命令
		fmt.Printf("→ %s\n", action)
		if err := e.commands.Exec(action); err != nil {
			fmt.Fprintf(os.Stderr, "执行失败: %v\n", err)
		}
	}

	return nil
}

// isDangerous 检查危险命令
func (e *AIEngine) isDangerous(cmd string) bool {
	dangerous := []string{"rm", "dd", "mkfs", "shutdown", "reboot", "> /dev/", "chmod 777"}
	for _, d := range dangerous {
		if strings.Contains(cmd, d) {
			return true
		}
	}
	return false
}
