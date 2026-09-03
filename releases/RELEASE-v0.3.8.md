## v0.3.8 · 2026-09-03 · Editorial 视觉改版

### 下载哪个文件？

| 你的系统 | 下载 |
|----------|------|
| **Windows 64 位** | `codebuddy-proxy-windows-x64.exe` |
| Linux x64 | `codebuddy-proxy-linux-amd64` |
| macOS Apple 芯片 | `codebuddy-proxy-darwin-arm64` |
| macOS Intel | `codebuddy-proxy-darwin-amd64` |

> Windows 请勿下载 `darwin-*` 或 `linux-amd64`。兼容旧名 `codebuddy-proxy-windows-amd64.exe`（与 x64 相同）。

校验：`SHA256SUMS.txt`

---

### 视觉

- 管理台 / OAuth 落地页 / GitHub Pages 产品页改为 **Editorial** 单色杂志风（`#F9F8F6` + `#1C1C1C`）
- Logo 单色重写；移除彩色光晕、噪点、Canvas 氛围层、Web Font CDN
- **仅展示层**：后端路由、账号池、OAuth、provider 行为不变

### 文档

- `design/` 规格改为 Editorial；`CHANGELOG` / 产品页下载表同步 Windows x64 推荐名

**说明**：预编译包需升级才能看到新管理台皮肤；纯静态产品页随仓库 `docs/` 更新到 GitHub Pages。
