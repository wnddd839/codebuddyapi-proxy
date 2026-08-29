# 更新日记

> CodeBuddy Proxy（Go · `protocol_direct`）发布说明。  
> 主仓库：[`wnddd839/proxy-codebuddy`](https://github.com/wnddd839/proxy-codebuddy)  
> 镜像：[`wnddd839/codebuddy-proxy`](https://github.com/wnddd839/codebuddy-proxy)

---

## 2026-08-29 · 审计硬化 + 体验收口

### 第二 / 第三梯队（代码卫生 + 测试）
- 删除 staticcheck U1000 死代码 4 处：`stainlessOS` / `mapOpenAIMessage` / `payloadBytes` / `Server.handle`。
- `NormalizeSite` 去掉恒为 global 的空串死分支。
- `/v1/models` 缓存 miss 用本地 singleflight 合并并发回源（无新依赖）。
- `server` 补 HTTP 层单测：API Key / 管理台免密与 BasicAuth / CSRF 同源接线。
- 重复工具收敛到 `strutil`：`Truncate` / `RandomHex` / `MaskSecret` / `Compact`。
- `oauth` / `billing` 补纯函数单测（URL、ShouldRefresh、JWT、用量汇总、通知码）。

### 二进制改走 GitHub Releases
- 仓库不再跟踪 `go-codebuddy/releases/codebuddy-proxy-*` 预编译包（约 30MB）。
- 二进制通过 `gh release` 发布到 GitHub Releases；Git 只保留源码与 `releases/README.md` / `.env.example`。
- 下载入口：https://github.com/wnddd839/proxy-codebuddy/releases/latest


本次把外部审计里的 **Batch A / B / C** 全部落地，并补齐此前已合入但未完整成文的体验修复。本地 `main` 相对远程新增提交：

| Commit | 摘要 |
|--------|------|
| `130a966` | 透出上游模型 `credits` 倍率；修复 UA 导致模型列表退化 |
| `c1ceb28` | 新增更新日记；管理台密码可留空（本地免密） |
| `c6f4f32` | BSD-3-Clause `LICENSE`、`.gitattributes`、清理本地垃圾 |
| `acf348d` | OAuth session 不再外泄可变指针；锁内鉴权 |
| `3f782f1` | 运行时配置改为 `atomic.Pointer[config.Config]` |
| `ad91e20` | 刷新 Windows amd64 二进制 |
| `a5bb211` | 模型缓存 / 采样透传 / 真实 `prompt_tokens` / CSRF / 去 Google Fonts |

### 1. 仓库卫生与许可证
- 补齐根目录 `LICENSE`（**BSD-3-Clause**），与 README badge 一致。
- 增加 `.gitattributes`：`*.go` 强制 LF，降低 Windows CRLF / gofmt 漂移。
- `.gitignore` 忽略本地垃圾：`.cline/`、`.omo/`、字面量 `~/`、`preview.html`、`*.exe~`。
- 清理误生成的 `~/`（含明文 OAuth session）与 IDE 临时目录。

### 2. OAuth session 并发安全
- **根因**：launch/callback 在锁外持有共享 `*OAuthSession`，与 `StartOAuth` 原地重置并发时 data race。
- `LiveOAuthSession()` 改为返回**值快照**，不再外泄可变指针。
- launch/callback 统一走锁内 `OAuthLaunchAuthorized(id, token)`。
- 增加 `go test -race` 覆盖的并发重置 / 鉴权测试；更新 `docs/operations/runbook.md`。

### 3. 运行时配置并发安全
- `gateway.Service` 用 go-modern `atomic.Pointer[config.Config]` 保存唯一运行时配置。
- 管理台改 API Key / 号池站点只更新这一份快照。
- `Server` 不再持有可写 `Cfg` 副本；`authorizeAPI` / `client-config` 一律读 `Svc.Config()`。
- 补 `-race` 测试（API Key + 号池切换并发读写）。

### 4. 模型列表缓存
- `/v1/models` **默认 60s TTL 缓存**，减轻频繁打 `/v3/config`。
- `?fresh=1` 或管理台「刷新模型」强制回源。
- 一键切换号池区域时自动清空缓存。

### 5. 模型倍率透出 + UA 修复
- 上游 `/v3/config` 自带 `credits`（如 `x0.00 credits` / `x0.29 credits`），归一化时不再丢弃。
- 管理台与 `/v1/models` 透出：`credits`、`creditMultiplier`、`free`。
- 可区分同名模型：`hy4-preview`（免费）vs `hy4-preview-x`（收费）等。
- User-Agent 改为 `CLI/<ver> CodeBuddy/<ver>`，避免裸 `CLI/<ver>` 触发上游 `12403 check ua` 导致列表退化成只有 `auto`。
- `PublicModelID` 改用 `strings.CutPrefix`，保留后缀原大小写。

### 6. Chat 采样参数透传
- OpenAI 兼容入口透传：
  - `temperature`
  - `top_p`
  - `max_tokens` / `max_completion_tokens`
- 未传字段不强制写入上游，保持默认行为。

### 7. 用量统计（`prompt_tokens`）
- 修复 `snapshot("")` 导致 `prompt_tokens` 恒为 1。
- 按真实 prompt 文本估算 token；若上游 SSE 带回 `usage`，**优先采用上游值**。

### 8. 管理台体验与安全
- `CODEBUDDY_PROXY_ADMIN_PASSWORD` **留空 = 本地免密打开管理台**。
- **`/v1` API Key 仍然保留**（`REQUIRE_API_KEY` / `API_KEY`），不再把空管理密码回填成 API Key。
- 去掉 Google Fonts，改系统字体栈（离线 / 内网更稳）。
- 管理台写操作增加 Origin / Referer **同源校验**（浏览器跨站 POST → 403；无头 curl 不受影响）。

### 9. 号池一键国内 / 国际切换（同周期已合入）
- 管理台一键切换 `domestic` / `global`。
- 切换后只使用当前区域账号；写入 `.env` 的 `CODEBUDDY_SITE` / `BASE_URL` / `INTERNET_ENVIRONMENT`。
- 账号自身 region 仍决定上游端点（国内 `copilot.tencent.com` / 国际 `www.codebuddy.ai`）。

### 10. 流式与路由硬化（同周期已合入）
- SSE：去掉过短总超时、写流加锁、提前打开 SSE。
- 客户端主动断开（`context canceled`）记为断开，不污染账号 `lastError`。
- keep-alive **5s**；响应头等待 **180s**。

### 本地验证清单
- `go test ./...` 通过；关键包带 `-race`。
- 管理台 `http://127.0.0.1:32126/direct-admin/` → 200（免密）。
- `/v1/models` + API Key → 200（约 28 模型，含 credits）。
- 跨站 `Origin: https://evil.example` POST 管理台 → **403**。
- 同源 Origin POST 号池切换 → **200**。

### 升级提示（二进制用户）
1. 替换 `releases/codebuddy-proxy-windows-amd64.exe`（或对应平台二进制）。
2. 保留现有 `.env`；API Key / 号池站点会继续生效。
3. 客户端超时的慢模型可优先用 `auto` / `deepseek-v4-flash`。
4. 需要强制刷新模型列表：`GET /v1/models?fresh=1`。

---

## 2026-08-28 · 首次启动 / 二进制体验

- 自动加载附近 `.env`（cwd / 可执行文件旁 / 上级目录）。
- 无 API Key 时首次启动自动生成并持久化到 `.env`，重启不丢。
- 默认 IDE 版本对齐 `2.117.2`。
- README / releases 补充上手说明。

---

## 历史里程碑（摘要）

- **Go 成为产品主线**：`protocol_direct` 网关；Pages：`main` + `/docs`。
- **主发布仓库**：`proxy-codebuddy`；`codebuddy-proxy` 作镜像。
- Node 版仅作本地历史备份，不再作为主发布路径。
