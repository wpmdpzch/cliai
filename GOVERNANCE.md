# ShellMind 项目治理

## 版本策略

### 三版本定位

| 版本 | 优先级 | 维护者 | 角色 |
|------|--------|--------|------|
| **shellmind (Go)** | P0 - 主推 | @wpmdpzch | 旗舰版本，最佳用户体验 |
| **shellmind-py (Python)** | P1 | @wpmdpzch | 开发者版，易扩展 |
| **shellmind-rs (Rust)** | P2 | @wpmdpzch | 极客版，性能优先 |

### 版本发布节奏

- **Go 版本**: 每月一个 minor 版本，每季度一个 major 版本
- **Python 版本**: 跟随 Go 版本同步发布
- **Rust 版本**: 每季度发布一次，与 Go major 版本同步

### 版本号规范

遵循 SemVer 2.0：

```
主版本.次版本.修订号
  v1.2.3

- MAJOR: 不兼容的 API 变更
- MINOR: 向后兼容的功能新增
- PATCH: 向后兼容的 bug 修复
```

### 功能同步策略

1. 新功能先在 Go 版本实现
2. 验证稳定后，同步到 Python 版本
3. Rust 版本保持核心功能同步，UI 特性可独立演进

---

## 代码规范

### Go 版本
- 遵循 [Effective Go](https://golang.org/doc/effective_go)
- 使用 `gofmt` 格式化
- 依赖管理: Go Modules
- 最低 Go 版本: 1.21

### Python 版本
- 遵循 PEP 8
- 类型注解: mypy strict mode
- 依赖管理: Poetry 或 uv
- 最低 Python 版本: 3.8

### Rust 版本
- 遵循 [Rust Style Guide](https://github.com/rust-dev-tools/fmt-rfcs)
- 使用 `cargo fmt`
- 依赖管理: Cargo
- 最低 Rust 版本: 1.70

---

## 分支策略

### Git Flow

```
main          # 稳定版本，始终可发布
├── develop   # 开发主分支
│   ├── feature/xxx   # 功能分支
│   ├── fix/xxx       # 修复分支
│   └── release/v1.x  # 发布分支
```

### 分支命名

```
feature/<功能名>
fix/<问题描述>
docs/<文档类型>
refactor/<重构内容>
```

---

## PR 流程

1. Fork 并创建 feature 分支
2. 遵循各语言代码规范
3. 提交 PR 到 `develop` 分支
4. 需要 1 个 reviewer 通过
5. CI 测试通过
6. Merge 后自动部署

---

## 决策机制

### RFC 流程

重大变更需要 Request for Comments：

1. 在 `docs/rfc/` 创建 RFC 文档
2. 描述问题和解决方案
3. 社区讨论 2 周
4. 核心团队最终决定

### 共识类型

- **Urgent**: 核心维护者直接决定
- **Normal**: 多数同意即可
- **Major**: 需要社区投票

---

## 社区

- Issue: GitHub Issues
- Discussion: GitHub Discussions
- 交流群: 待定
