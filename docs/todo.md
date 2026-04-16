# CLI-AI v0.1 Todo List

> 项目进度跟踪 | 更新: 2026-04-16

---

## 编译产物
- [x] `cliai-go/cliai` 二进制文件 (11MB)

---

## 功能待完成

### P0 - 必须完成
- [ ] **Tab 键检测**：TUI 中实现真正的 Tab 键切换模式
- [ ] **命令执行**：TUI 中输入的命令实际调用 pkgcmd 执行
- [ ] **AI 集成**：自然语言解析，调用 LLM 生成命令

### P1 - 应该完成
- [ ] **curl 增强**：支持 POST、Header、Cookie 等
- [ ] **jq 增强**：支持数组索引、[n]、@unique 等
- [ ] **base64 文件支持**：支持从文件读取
- [ ] **帮助文档**：内置命令帮助信息

### P2 - 尽量完成
- [ ] **其他 P0 命令**：grep, cat, ls, ps, df, free, head, tail
- [ ] **包管理**：install 命令
- [ ] **配置文件**：~/.cliai/config.yaml 读取

---

## 问题记录

### Tab 键检测
- 当前 `bufio.Reader` 无法检测 Tab 键
- 需要使用原始模式 `stty cbreak` + 单字节读取
- 方案：改用 `syscall.Read` 直接读取

### TUI 显示重叠
- 当前输入提示和边框有重叠
- 方案：调整光标位置或使用完整重绘

### 网络限制
- Go 模块下载超时（golang.org, proxy.golang.org）
- 解决方案：使用 `GOPROXY=https://goproxy.cn,direct`
- 部分库版本过高（如 bubbletea 需要 go 1.24）

---

## 已完成功能
- [x] Go 1.22.4 环境安装
- [x] 项目框架搭建
- [x] 内置命令：base64, curl, jq（Go 原生实现）
- [x] TUI 窗口模式
- [x] REPL 交互模式
- [x] 命令行子命令架构

---

## 下一步优先级
1. Tab 键检测（TUI 体验关键）
2. 命令执行（让 TUI 有实际功能）
3. AI 集成（核心价值）
