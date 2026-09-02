## v0.3.7 · 2026-08-31 · 号池冷却与规范落地

### 下载哪个文件？

| 你的系统 | 下载 |
|----------|------|
| **Windows 64 位** | `codebuddy-proxy-windows-x64.exe` |
| Linux x64 | `codebuddy-proxy-linux-amd64` |
| macOS Apple 芯片 | `codebuddy-proxy-darwin-arm64` |
| macOS Intel | `codebuddy-proxy-darwin-amd64` |

> Windows 请勿下载 `darwin-*` 或 `linux-amd64`。

校验：`SHA256SUMS.txt`

---

### 号池冷却

- 上游失败写入账号 `cooldownUntil`（11140/11128/11101/11102 → 5min，429 → 2min，502/503/504 → 30s）
- **全冷却降级**：同区域全部账号冷却时，选最快恢复者继续服务
- **11140** 不再触发换号重试

### OAuth / 账号

- Upsert 去重改为 `(userId + site)`，国内/国际 OAuth 可并存
- 旧版已合并账号需对缺失区域 **重新 OAuth**

### 路由

- 新增 **`GET /model/info`** 别名（OpenCode `modelInfoEndpoint: "/model/info"`）

### 工程

- pre-commit：gofmt 检查 + 包测试
- Pi 系统提示词：`.agents/pi-system-prompt-go.md`

**linux-amd64 SHA256:** `d4695e8a89985376911f830aee17a8b73000c3edc29690080350ec6b7a72834a`
