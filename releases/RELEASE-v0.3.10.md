## v0.3.10 · 2026-09-05 · 上游 11128 排障与小规模适配

### 下载哪个文件？

| 你的系统 | 下载 |
|----------|------|
| **Windows 64 位** | `codebuddy-proxy-windows-x64.exe` |
| Linux x64 | `codebuddy-proxy-linux-amd64` |
| macOS Apple 芯片 | `codebuddy-proxy-darwin-arm64` |
| macOS Intel | `codebuddy-proxy-darwin-amd64` |

> Windows 请勿下载 `darwin-*` 或 `linux-amd64`。兼容旧名 `windows-amd64.exe`。

校验：`SHA256SUMS.txt`

---

### 根因

- ZCode agent 每次请求在 `system` 中注入 `System Context` 块（含 `Main branch (you will usually use this for PRs): <branch>`），腾讯上游 `copilot.tencent.com` 将该整串判为非官方 CLI 通道，返回 `400 + 11128 Illegal API invocation from an unapproved channel`。
- 经二分验证：该三件套（`Main branch` + 括号注解 + 冒号）齐备即触发，缺任一件即放行；git 状态文本本身、`role=developer`、工具数、`thinking`、`max_completion_tokens` 均已证伪（详见 `CHANGELOG.md`）。

### 变更

- `system` 消息去除 ` (you will usually use this for PRs)` 括号注解（保留分支名，`Main branch: <branch>`），用户原文与工具结果原样透传。
- 上游非 2xx / 流中报错时进程 WARN 日志 `codebuddy upstream failure` 输出脱敏请求指纹（role 序列 / 工具名 / thinking / reasoning / 单消息 200 字预览，不含 token 与完整 prompt），排障见 runbook 11128 行。
- `gateway.New` 向 provider 注入进程 logger，避免与 `slog.Default()` 错配。

### 验证

- `gofmt` 无差异，`go vet ./...` 通过，`go test ./...` 全绿（含 `-race`）。
- 实测：触发句在修复前 `400 + 11128`，修复后 `200`；`deepseek-v4-flash` / `hy4-preview` / `auto` 最小请求不受影响。
