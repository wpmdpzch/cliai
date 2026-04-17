# CLI-AI v0.1 Todo List

> 项目进度跟踪 | 更新: 2026-04-17

---

## 编译产物
- [x] `cliai-go/cliai` 二进制文件 (11MB)

---

## 功能待完成

### P0 - 必须完成
- [x] **Tab 键检测**：TUI 中实现真正的 Tab 键切换模式
- [x] **命令执行**：TUI 中输入的命令实际调用 pkgcmd 执行
- [x] **AI 集成**：自然语言解析，调用 LLM 生成命令

### P1 - 应该完成
- [x] **curl 增强**：支持 POST、Header、Cookie 等
- [x] **jq 增强**：支持数组索引、[n]、@unique 等
- [x] **base64 文件支持**：支持从文件读取
- [x] **帮助文档**：内置命令帮助信息

### P2 - 尽量完成
- [x] **其他 P0 命令**：grep, cat, ls, ps, df, free, head, tail
- [ ] **包管理**：install 命令
- [x] **配置文件**：~/.cliai/config.yaml 读取

---

## 问题记录

### 已修复 (2026-04-17)
- [x] **REPL panic**: nil pointer dereference - 修复：添加 config 加载逻辑
- [x] **stdin 命令失效**: grep/wc/sed/awk - 修复：execSystemCmd 添加 stdin 检测
- [x] **jq 数组索引**: 不支持 [0] - 修复：实现 strconv.Atoi 解析
- [x] **md5/sha256 命令名**: 系统命令是 md5sum/sha256sum - 修复：更新命令注册名
- [x] **curl POST 数据**: -d 未传递 body - 修复：strings.NewReader(data)
- [x] **错误信息不清晰**: stderr 混入 stdout - 修复：分离 stdout/stderr 输出
- [x] **帮助信息简陋**: 无参数时无提示 - 修复：无参数自动显示帮助
- [x] **commands 输出单调**: 无分类无图标 - 修复：emoji 分类 + 彩色输出
- [x] **二进制内容乱码**: cat 等命令输出二进制时终端混乱 - 修复：二进制检测

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
- [x] stdin 命令支持（grep, wc, sed, awk）
- [x] jq 数组索引支持
- [x] md5sum/sha256sum 命令修正
- [x] curl POST body 支持
- [x] 配置文件加载
- [x] 彩色输出（错误/警告/提示）
- [x] emoji 分类显示
- [x] 二进制内容检测
- [x] 欢迎信息边框
- [x] 帮助信息增强

---

## 下一步优先级
1. Tab 键检测（REPL 体验关键）
2. 包管理功能
3. AI 实际集成测试
4. TUI 显示优化
