# CLI-AI

> **让命令行会思考**
>
> A natural language CLI that speaks your language and speaks to your system.

**自然语言 / 原生命令 → 统一执行**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.8+-3776AB.svg)](https://python.org)
[![Rust Version](https://img.shields.io/badge/Rust-1.70+-CE422B.svg)](https://rust-lang.org)

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

> 🎯 **自带命令包**：内置 `curl`、`jq`、`grep`、`sed`、`awk` 等常用工具，不依赖系统环境，跨平台开箱即用。

---

## 能做什么？

```bash
# 排查服务器为什么慢
$ cliai "帮我检查服务器性能，列出最吃资源的进程"
→ 自动生成并执行性能分析命令

# 批量处理文件
$ cliai "把所有 json 文件转成 csv"
→ 内置 jq/awk，自动识别并转换

# 跨平台兼容
$ cliai "帮我监控 80 端口的连接数"
→ Linux/macOS/Windows 自动适配

# 网络调试
$ cliai "发送一个 GET 请求到 example.com"
→ 内置 curl，无需系统安装
```

---

## 三版本架构

| 版本 | 优势 | 安装方式 | 目标用户 |
|------|------|---------|---------|
| **Go** | 编译单一二进制，无依赖，跨平台最强 | `curl -fsSL https://get.cliai.dev \| bash` | 通用用户 / 服务器运维 |
| **Python** | 开发快，生态丰富，易魔改 | `pip install cliai` | 开发者 / 脚本爱好者 |
| **Rust** | 性能最优，内存安全 | `cargo install cliai` | 性能党 / 极客 |

---

## 快速上手

### Go 版本（推荐）

```bash
# 一键安装（Linux/macOS）
curl -fsSL https://get.cliai.dev | bash

# Windows
winget install cliai

# 或手动安装
go install github.com/wpmdpzch/cliai/cliai-go@latest
```

### Python 版本

```bash
pip install cliai
# 或
git clone https://github.com/wpmdpzch/cliai && cd cliai-py && pip install -e .
```

### Rust 版本

```bash
cargo install cliai
# 或
git clone https://github.com/wpmdpzch/cliai && cd cliai-rs && cargo build --release
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

> 💡 如果系统已安装对应命令，优先使用系统命令（版本更新、更多选项）。

---

## 配置

首次运行会自动创建配置，也支持手动创建：

```yaml
# ~/.cliai/config.yaml
ai:
  provider: "openai"           # openai / claude / gemini / ollama / azure
  api_key: "sk-xxxxx"
  base_url: "https://api.openai.com/v1"
  model: "gpt-4o-mini"
  temperature: 0.7
  max_tokens: 2048

exec:
  auto_exec: true              # false 则只显示命令
  confirm_dangerous: true      # 危险命令确认
  timeout: 30                  # 秒

ui:
  mode_indicator: true         # 显示 CLI/PLAN/BUILD 模式
  default_mode: "cli"          # 默认模式

tools:
  # 内置命令包
  builtin: ["network", "text", "file", "system", "encoding"]
  # 扩展工具
  enabled: ["docker", "git"]
```

---

## 项目结构

```
cliai/
├── cliai-go/                   # Go 版本（主推）
│   ├── cmd/cliai/             # 入口
│   ├── core/                  # 核心解析引擎
│   ├── builtin/               # 内置命令包
│   ├── ui/                    # OpenCode 窗口模式
│   ├── tools/                 # 扩展工具链
│   └── config/                # 配置管理
│
├── cliai-py/                   # Python 版本
│   └── cliai/                 # 包
│
├── cliai-rs/                   # Rust 版本
│   └── src/                   # 源码
│
├── docs/                       # 文档
├── ROADMAP.md                  # 路线图
└── GOVERNANCE.md               # 项目治理
```

---

## 参与贡献

1. Fork → Feature Branch → PR
2. 遵循三版本代码规范
3. 所有 PR 需要通过测试

---

## License

MIT © wpmdpzch
