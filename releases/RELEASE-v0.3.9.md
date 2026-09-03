## v0.3.9 · 2026-09-03 · 管理台章节切页

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

### 变更

- 管理台改为四章节切页（概览 / 账号池与 OAuth / 客户端接入 / 模型与快照）
- `#codebuddy` / `#client-config` hash 自动切到对应章节
- 概览文案与号池冷却策略对齐
- **后端逻辑未改**；仅 `internal/admin` 展示层 + `design/01` 规格

预编译用户需升级到本版本才能看到切页管理台。
