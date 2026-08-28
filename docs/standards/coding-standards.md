# Coding Standards

本仓库 Go 编码规范 = **JetBrains Modern Go Guidelines** + 项目约束。

## 权威来源

1. 项目 skill：`.agents/skills/use-modern-go`  
   上游：https://github.com/JetBrains/go-modern-guidelines  
   （JetBrains / IntelliJ IDEA / GoLand 官方 Agent Skill）
2. 摘录文档：`docs/standards/modern-go-guidelines.md`
3. 原文归档：`docs/standards/vendor/`

写 Go 前先执行：

```powershell
.\.agents\skills\use-modern-go\scripts\run-tool.ps1 list --go-version 1.26
```

## 本项目强制约定

### 语言与工具链

- Go version 以 `go.mod` 为准（当前 1.26.x）
- `gofmt` 必须通过
- 优先标准库；新增第三方依赖需写进 PR/文档说明

### 现代语法优先（节选）

| 场景 | 使用 |
|------|------|
| 首个非空值 | `cmp.Or` / `strutil.First` |
| 前缀判断并切片 | `strings.CutPrefix` / `CutSuffix` |
| 最值 | 内建 `min` / `max` |
| 成员判断 | `slices.Contains` |
| 路由 | Go 1.22 method-aware `ServeMux` |
| 错误类型匹配 | `errors.AsType[T]`（Go 1.26） |
| 指针字面量 | `new(value)`（Go 1.26） |
| JSON 零值省略 | `omitzero`（适合 bool/number/struct/time） |

### 工程结构

- 业务代码只放 `internal/`
- 入口只放 `cmd/`
- 禁止把协议细节塞进 `main`
- 禁止恢复 Node 时代的多传输默认路径

### 安全

- 禁止日志输出完整 token
- 账号文件权限保持私有（0600）
- Admin 与 API 鉴权分离

### 测试

至少覆盖：

- accounts 默认 enabled / 轮询
- config 默认值
- provider endpoint / SSE / models normalize

```bash
go test ./...
```

### 文档同步

变更下列内容时必须更新 `docs/`：

- 路由
- 环境变量
- 包边界
- OAuth / 账号池行为
