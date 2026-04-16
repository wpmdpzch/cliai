# ShellMind

> **让 Shell 会思考**
>
> A natural language CLI that speaks your language and speaks to your system.

**自然语言 → 系统命令 → 自动执行**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.8+-3776AB.svg)](https://python.org)
[![Rust Version](https://img.shields.io/badge/Rust-1.70+-CE422B.svg)](https://rust-lang.org)

---

## 能做什么？

```bash
# 排查服务器为什么慢
$ sm "帮我检查服务器性能，列出最吃资源的进程"
→ htop, top, ps aux 等命令已生成并执行

# 批量处理文件
$ sm "把所有 json 文件转成 csv"
→ 自动识别文件，生成转换命令并执行

# 跨平台兼容
$ sm "帮我监控 80 端口的连接数"
→ Linux: netstat/ss | macOS: lsof | Windows: netstat
```

---

## 三版本架构

| 版本 | 优势 | 安装方式 | 目标用户 |
|------|------|---------|---------|
| **Go** | 编译单一二进制，无依赖，跨平台最强 | `curl -fsSL https://get.shellmind.dev \| bash` | 通用用户 / 服务器运维 |
| **Python** | 开发快，生态丰富，易魔改 | `pip install shellmind` | 开发者 / 脚本爱好者 |
| **Rust** | 性能最优，内存安全 | `cargo install shellmind` | 性能党 / 极客 |

---

## 快速上手

### Go 版本（推荐）

```bash
# 一键安装（Linux/macOS）
curl -fsSL https://get.shellmind.dev | bash

# Windows
winget install shellmind

# 或手动安装
go install github.com/wpmdpzch/cliai/shellmind@latest
```

### Python 版本

```bash
pip install shellmind
# 或
git clone https://github.com/wpmdpzch/cliai && cd shellmind-py && pip install -e .
```

### Rust 版本

```bash
cargo install shellmind
# 或
git clone https://github.com/wpmdpzch/cliai && cd shellmind-rs && cargo build --release
```

### Docker

```bash
docker run -it wpmdpzch/cliai:latest
```

---

## 配置

首次运行会自动创建配置，也支持手动创建：

```yaml
# ~/.shellmind/config.yaml
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

tools:
  enabled: ["system", "network", "process", "file", "git", "docker"]
```

---

## 项目结构

```
cliai/
├── shellmind/                  # Go 版本（主推）
│   ├── cmd/shellmind/          # 入口
│   ├── core/                   # 核心解析引擎
│   ├── tools/                  # 工具链
│   └── config/                 # 配置管理
│
├── shellmind-py/               # Python 版本
│   └── shellmind/              # 包
│
├── shellmind-rs/               # Rust 版本
│   └── src/                    # 源码
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

MIT © wpm-flying-nest
