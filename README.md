# CLI-AI

> **让命令行会思考**
>
> A natural language CLI that speaks your language and speaks to your system.

**自然语言 / 原生命令 → 统一执行**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org)

---

## 交互模式

采用 **OpenCode 窗口模式**，Tab 键切换三种模式：

| 模式 | 说明 | 能力 |
|------|------|------|
| **CLI** | 命令行直接执行 | 原生命令、内置命令包 |
| **PLAN** | 预执行（AI） | readonly，只读操作 |
| **BUILD** | 直接操作（AI） | write，可写文件 |

```bash
# Tab 键循环切换
[CLI] → [PLAN] → [BUILD] → [CLI] ...

# CLI 模式：直接执行命令
$ cliai "curl https://api.test.com/ping"

# PLAN 模式：AI 理解，但只读预览
[PLAN]$ cliai "帮我分析这个日志文件"
→ AI 解析并展示分析结果（不执行危险操作）

# BUILD 模式：AI 直接操作
[BUILD]$ cliai "把这个目录下的所有 .js 文件压缩成 one.js"
→ AI 执行写操作：读取、合并、写入
```

---

## 核心理念

**两个入口，一个目标**

```bash
# 自然语言 → AI 解析 → 执行
$ cliai "帮我检查服务器性能"
→ AI 理解意图，生成并执行命令

# 原生命令 → 直接执行（无需系统安装）
$ cliai "curl https://api.test.com/ping"
$ cliai "jq '.data' response.json"
$ cliai "grep -r 'error' /var/log"
```

> 🎯 **自带命令包**：内置 `curl`、`jq`、`grep`、`sed`、`awk` 等常用工具，优先调用，跨平台开箱即用。

---

## 能做什么？

```bash
# 排查服务器为什么慢
$ cliai "帮我检查服务器性能，列出最吃资源的进程"
→ 自动生成并执行性能分析命令

# 批量处理文件
$ cliai "把所有 json 文件转成 csv"
→ 内置 jq/awk，自动识别并转换

# 网络调试
$ cliai "发送一个 GET 请求到 example.com"
→ 内置 curl，无需系统安装

# 自然语言安装命令包
$ cliai "帮我添加一个图片处理的包"
→ AI 生成并安装命令脚本
```

---

## 版本策略

> ⚠️ **v0.1 仅开发 Go 版本**。Python/Rust 版本暂停。

| 版本 | 状态 |
|------|------|
| **Go** | ✅ 开发中 |
| **Python** | ⏸️ 暂停 |
| **Rust** | ⏸️ 暂停 |

> ⚠️ **v0.1 仅支持 Linux 系统**。macOS/Windows 后续版本支持。

---

## 快速上手

### 下载 Release（推荐）

**二进制方式：**

```bash
# x86_64 Linux
curl -L https://github.com/wpmdpzch/cliai/releases/latest/download/cliai-linux-amd64 -o cliai
chmod +x cliai
sudo mv cliai /usr/local/bin/

# 验证
cliai --version
```

**源代码方式：**

```bash
git clone https://github.com/wpmdpzch/cliai.git
cd cliai/cliai-go
go build -o cliai ./cmd/cliai
./cliai --version
```

### Docker

```bash
docker run -it wpmdpzch/cliai:latest
```

---

## 内置命令包

CLI-AI 自带常用工具，不依赖系统：

| 类别 | 内置命令 |
|------|---------|
| **网络** | curl, wget, ping, nc |
| **文本** | jq, grep, sed, awk, cut, sort, uniq |
| **文件** | ls, cat, head, tail, wc, diff |
| **系统** | ps, top, df, du, free |
| **编码** | base64, md5, sha256 |

> ⚠️ **自带的优先调用，系统命令兜底**。

### 命令调用优先级

1. **自带命令包**（优先）
2. **系统命令**（兜底，warn 如果不存在）

### 命令来源

- **核心命令**：Go 原生实现（curl, jq 等）
- **系统命令**：调用系统，无自带时兜底

### 添加新命令包

**方式 1：自然语言添加**

```bash
$ cliai "帮我添加一个图片处理的包"
→ AI 理解需求，生成命令脚本
→ 用户选择安装方式
```

**方式 2：包管理器安装**

```bash
$ cliai "install image"
$ cliai "install json"
$ cliai "install git"
```

**方式 3：手动安装**

```bash
# 放入用户目录
~/.cliai/packages/image/resize
~/.cliai/packages/image/convert

# 或安装到系统
/usr/local/cliai/packages/image/
```

---

## 包管理

### 包结构

每个包一个目录：

```
image/                          # 包名
├── command.json               # 元数据（AI 读取）
├── resize                     # 可执行脚本
└── convert
```

**command.json 示例：**

```json
{
  "name": "image",
  "description": "图片处理工具集",
  "version": "1.0.0",
  "commands": [
    {
      "name": "resize",
      "usage": "image resize <file> <width> <height>",
      "description": "调整图片尺寸",
      "example": "image resize photo.jpg 800 600",
      "dangerous": false
    }
  ]
}
```

### 包存放位置

```yaml
# config.yaml
packages:
  # 用户目录（默认）
  local: "~/.cliai/packages"
  # 系统目录（可选）
  system: "/usr/local/cliai/packages"
```

### 命令冲突处理

同名命令存在时，逐一尝试执行，以能实现功能为准。

### 包来源

- **官方包**：GitHub Releases / 官方文档站
- **AI 生成**：自然语言描述需求，AI 生成
- **社区贡献**：GitHub PR

### 命令描述

- 开发者编写命令时必须写好描述
- AI 从 `command.json` 读取
- 可从官方文档查询补充

---

## 配置

首次运行会自动创建配置，也支持手动创建：

```yaml
# ~/.cliai/config.yaml
ai:
  provider: "openai"           # 通用配置，任意 provider
  api_key: "sk-xxxxx"
  base_url: "https://api.openai.com/v1"
  model: "gpt-4o-mini"
  temperature: 0.7
  max_tokens: 2048

exec:
  auto_exec: false             # true 直接执行，false 只显示命令
  confirm_dangerous: true      # 危险命令确认
  timeout: 30                  # 秒

ui:
  mode_indicator: true         # 显示 CLI/PLAN/BUILD 模式
  default_mode: "cli"          # 默认模式

packages:
  local: "~/.cliai/packages"   # 用户包目录
  system: "/usr/local/cliai/packages"  # 系统包目录
```

---

## 项目结构

```
cliai/
├── cliai-go/                   # Go 版本（主推）
│   ├── cmd/cliai/             # 入口
│   ├── core/                  # 核心解析引擎
│   ├── builtin/               # 内置命令包（Go 实现）
│   │   └── commands.json     # 内置命令清单
│   ├── pkg/                   # 包管理
│   │   ├── manager.go         # 包管理器
│   │   └── scanner.go         # 命令扫描器
│   ├── ui/                    # OpenCode 窗口模式
│   └── config/                # 配置管理
│
├── docs/                       # 文档
├── ROADMAP.md                  # 路线图
├── DESIGN.md                   # 设计文档
└── GOVERNANCE.md               # 项目治理
```

---

## 参与贡献

1. Fork → Feature Branch → PR
2. 遵循 Go 代码规范（Effective Go）
3. 所有 PR 需要通过测试

---

## License

MIT © wpmdpzch
