# 更新日记

## 2026-08-29

### OAuth session 并发安全
- `LiveOAuthSession()` 改为返回值快照（不再外泄可变指针）。
- launch/callback 统一走锁内 `OAuthLaunchAuthorized(id, token)`。
- 增加 `go test -race` 覆盖的并发重置/鉴权测试；更新 runbook。

### 仓库卫生 / 许可证
- 补齐根目录 `LICENSE`（BSD-3-Clause），与 README badge 一致。
- 增加 `.gitattributes`（`*.go` 强制 LF），避免 gofmt/CRLF 漂移。
- `.gitignore` 忽略本地垃圾：`.cline/`、`.omo/`、字面量 `~/`、`preview.html`、`*.exe~`。
- 清理本地误生成的 `~/`（含明文 OAuth session）与 IDE 临时目录。

### 模型倍率透出 + 模型列表修复
- 上游 `/v3/config` 本身就带 `credits` 倍率（如 `x0.00 credits` / `x0.29 credits`），之前反代归一化时丢掉了。
- 现在管理台与 `/v1/models` 都会透出：
  - `credits`
  - `creditMultiplier`
  - `free`
- 可区分同名模型：
  - `hy4-preview` = 免费（x0.00）
  - `hy4-preview-x` = 收费（x0.29）
  - `hy3` / `hy3-x` 同理
- 修复 User-Agent：裸 `CLI/<ver>` 会触发上游 `12403 check ua`，导致模型列表退化成只有 `auto`。现改为 `CLI/<ver> CodeBuddy/<ver>`。

### 管理台免密（本地更省事）
- `CODEBUDDY_PROXY_ADMIN_PASSWORD` 留空 = 管理台免密打开。
- **API Key 仍然保留**（`CODEBUDDY_PROXY_REQUIRE_API_KEY` / `CODEBUDDY_PROXY_API_KEY`）。
- 不再把空的管理密码自动回填成 API Key。

### 号池一键国内 / 国际切换
- 管理台账号池支持一键切换 `domestic` / `global`。
- 切换后只使用当前区域账号发起请求。
- 写入 `.env` 的 `CODEBUDDY_SITE` / `CODEBUDDY_BASE_URL` / `CODEBUDDY_INTERNET_ENVIRONMENT`，重启仍生效。

### 流式与路由硬化（同周期）
- SSE 更稳：去掉过短总超时、写流加锁、提前打开 SSE。
- 客户端主动断开（`context canceled`）记为断开，不再污染账号 `lastError`。
- keep-alive 5s；响应头等待放宽到 180s。
- 按账号自身 region 路由上游（国内 `copilot.tencent.com` / 国际 `www.codebuddy.ai`）。

---

## 2026-08-28

### 首次启动 / 二进制体验
- 自动加载附近 `.env`（cwd / 可执行文件旁 / 上级目录）。
- 无 API Key 时首次启动自动生成并持久化到 `.env`。
- 默认 IDE 版本对齐 `2.117.2`。
