# 更新日记

> CodeBuddy Proxy（Go · `protocol_direct`）发布说明。  
> 主仓库：[`wnddd839/codebuddyapi-proxy`](https://github.com/wnddd839/codebuddyapi-proxy)  
> 镜像：[`wnddd839/codebuddy-proxy`](https://github.com/wnddd839/codebuddy-proxy)

---

## v0.3.6 · 2026-08-31 · OpenCode / 号池 / 国际模型列表

### OpenCode 思考档位透传

- **根因**：OpenCode `@ai-sdk/openai-compatible` 需要 models.dev 形态字段（`reasoning: true` 布尔、`reasoning_options`、`variants`），此前 `/v1/models` 把上游 `reasoning` 对象直接透出，客户端无法识别思考档位 UI。
- **`GET /v1/models` 扩展**：支持思考的模型返回 `reasoning: true`、`reasoning_config`（上游元数据）、`reasoning_options`、`interleaved`、`variants`（按 `supportedEfforts` 生成，如 `glm-5.3-flash` 为 `low/high/max`）。
- **`GET /v1/model/info`**：LiteLLM 格式元数据，供 `opencode-models-discovery` 配置 `modelInfoFormat: "litellm"` 自动注入思考档位（provider 级一次配置，无需 per-model 手配）。
- **管理台**：不增加「复制 opencode.json」入口；透传 + discovery 插件即可。

### 号池 failover

- **429 / 502 / 503 / 504 / rate limit** 等可恢复上游故障：标记当前账号 `failedRequests` + `lastError`，自动换下一个号重试（最多 3 层，轮询选号 + `ExcludeIDs` 跳过已试账号）。
- **仍不换号**：`11128`（渠道限制）、`11101` / `11102`（请求/模型形态）、`401`/`403`（先 OAuth refresh 同一账号）。

### 国际站模型列表（hy4 缺失）

- **根因**：国际号池只打 `www.codebuddy.ai/v3/config`（约 35 个模型，含 Gemini/GPT，**无 hy4**）；`copilot.tencent.com/v3/config` 用同一国际 token 可读且含 `hy4-preview` 等。
- **修复**：国际站合并两个 `/v3/config` 源（去重合并）；国内站仍只走 `copilot.tencent.com`。聊天仍按账号区域走 `codebuddy.ai` / `copilot.tencent.com` 对应 chat 端点。

### 二进制启动体验

- 启动时输出简短中文提示，**引导打开管理台**完成 OAuth 与客户端接入配置。
- 默认日志级别 `warn`（排障设 `CODEBUDDY_PROXY_LOG_LEVEL=info|debug`）；**不再在日志中打印明文 API Key**。
- 号池切换日志降为 debug。

### 工具

- `go-codebuddy/scripts/probe-v3-models.js`：对比各上游 `/v3/config` 模型列表（排障用）。

### Release 资产命名

- **Windows 64 位请下载 `codebuddy-proxy-windows-x64.exe`**（与 `windows-amd64.exe` 相同，后者为兼容旧名）。
- 勿在 Windows PC 上下载 `darwin-arm64` / `darwin-amd64`（macOS 专用）。

下载：https://github.com/wnddd839/codebuddyapi-proxy/releases/tag/v0.3.6

---

## v0.3.5 · 2026-08-31 · 思考档位透传

- **Chat 入站**：解析 `reasoning_effort` / `reasoningEffort` / `reasoning` / `thinking`，写入上游 body。
- **上游对齐**：实测仅 `reasoning.effort` 触发 `reasoning_content`；网关将 `reasoning_effort` 映射为 `reasoning: { effort }`（同时保留 `reasoning_effort` 兼容字段）。
- **模型列表**：`/v1/models` 透传 `supportsReasoning` / `onlyReasoning` / `reasoning`（含 `supportedEfforts`、`defaultEffort`）。
- **响应**：流式 delta 与 non-stream `message.reasoning_content` 回传思考内容；无思考时 `omitempty` 零回归。

下载：https://github.com/wnddd839/codebuddyapi-proxy/releases/tag/v0.3.5

---

## v0.3.4 · 2026-08-30 · 出站 cache 字段别名补齐

- **问题**：管理台能看到 `totalCachedTokens`，但 CCSwitch 等下游缓存率仍极低。
- **核实**：流式收尾 `usage` chunk（空 choices）本身已在发送；CodeBuddy 上游多用 DeepSeek 风格 `prompt_cache_*` 字段。
- **修复**：`UsageFromProvider` 在命中时同时写出 `prompt_tokens_details.cached_tokens` + `prompt_cache_hit_tokens` + `cache_read_input_tokens`，并推导 `prompt_cache_miss_tokens`；`ParseUsage` 兼容 `cached_tokens` / `input_tokens_details`。
- **流式**：usage chunk 写出后显式 `Flush()`，降低下游提前断开丢统计的概率。

下载：https://github.com/wnddd839/codebuddyapi-proxy/releases/tag/v0.3.4

---

## v0.3.3 · 2026-08-29 · 回滚出站混合序列化（热修复）

- **回滚**：删除 `MarshalStreamChunk` / `SSEMarshaler` 手写骨架路径，出站 `StreamChunk` 恢复 `json.Marshal`。
- **原因**：typed struct 的 `json.Marshal` 已是 **2 allocs / ~288B**；混合路径 **9 allocs / ~552B**，实测为负优化（上轮基线误用 `map[string]any`）。
- **保留**：v0.3.1 SSE 缓冲聚合、v0.3.0 上游 typed 解析、中文 token 估算等有效优化不受影响。

下载：https://github.com/wnddd839/codebuddyapi-proxy/releases/tag/v0.3.3

---
