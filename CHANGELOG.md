# 更新日记

> CodeBuddy Proxy（Go · `protocol_direct`）发布说明。  
> 主仓库：[`wnddd839/codebuddyapi-proxy`](https://github.com/wnddd839/codebuddyapi-proxy)  
> 镜像：[`wnddd839/codebuddy-proxy`](https://github.com/wnddd839/codebuddy-proxy)

---

## v0.3.0 · 2026-08-29 · 速度 / 体验 / 内存占用全面收口

本版是 v0.2.0 之后的性能与体验大版本：账号池、流式写出、token 统计、管理台与工程卫生一起落地。  
目标很明确——**更省内存、更快响应、更少踩坑**，二进制用户开箱即用。

### 亮点一览

| 维度 | 变化 |
|------|------|
| **速度** | 账号池热路径从「每请求同步写盘」→ 内存权威 + 250ms 异步刷盘；SSE 写出加 16KB 缓冲；正则预编译；SSE typed 解析 |
| **体验** | token/缓存透传给 CCSwitch；管理台总 Tokens；流式更稳；免密管理台；暗色 instrument UI；产品页 / logo / README 定稿 |
| **内存 / 占用** | 进程 RSS 约 **5–7MB**；去掉 Git 内嵌二进制（仓库体积大减）；刷盘合并后磁盘抖动消失；panic 不会拖垮进程 |
| **正确性** | 中文 token 不再按字节低估；`usage_source` 区分上游/估算；重试深度上限；OAuth/配置竞态修完 |

下载：https://github.com/wnddd839/codebuddyapi-proxy/releases/tag/v0.3.0

---

### 1. 速度：热路径砍掉「无用等待」

#### 账号池：内存权威 + 异步刷盘（最大收益之一）
- `Select` / `MarkResult` 只改内存并标记 dirty，**250ms 防抖合并写盘**。
- 以前每个请求约 4 次同步磁盘 IO（全局锁 + `MarshalIndent` + rename），并发实测约 **~38 ops/s**。
- 现在并发 Select+Mark 可达 **数十万 ops/s**；凭据变更（Upsert/Delete/Replace/SetEnabled）仍同步落盘，保证不丢号。
- `json.MarshalIndent` → `json.Marshal`；`MkdirAll` 只做一次；`Close`/`Shutdown` 强制 flush。
- `persistMu` 与业务 `mu` 分离；`flushLoop` 正确 drain timer，避免 goroutine 泄漏。

#### 流式 SSE：缓冲写出 + 少分配
- 新增 `httputil.SSEStream`：16KB `bufio.Writer` + `fmt.Appendf` 组帧，再按帧 Flush。
- 避免「每个 chunk 一次 syscall / 字符串拼接」；Windows 小包写开销更高，收益更明显。
- 上游 SSE 热路径优先 **typed OpenAI chunk** 解析，失败才回退 `map[string]any`（分配明显下降）。
- 模型名解析 / 鉴权与换号判断正则改为**包级预编译**（热路径约 40×）。

#### 请求收尾不再白烧 CPU
- 上游已返回 `usage` 时，**不再**无条件 `estimatePromptText`（长上下文可省百微秒级）。
- `CompleteFromPool` 增加 `RetryDepth`（max=3），换号 + OAuth 刷新不会指数展开。

---

### 2. 体验：客户端与管理台「看得懂、用得顺」

#### Token / 缓存命中透传
- 流式收尾补发 OpenAI `include_usage` 风格 usage chunk（空 choices + usage），CCSwitch 等能统计总 token / 缓存。
- 透传 `prompt_tokens_details.cached_tokens`，并兼容 Anthropic / DeepSeek 别名。
- 管理台 Stats 累计 `totalTokens` / `totalCachedTokens`，首页新增「总 Tokens」。
- `usage_source=upstream|estimated`：客户端能区分真实值与本地估算。
- 中文 token 估算改为按 rune/CJK，修复「字节/4」系统性低估 50%+。

#### 流式稳定性
- 去掉 `http.Client` 总 Timeout（长 agent / 慢模型不会被中途掐断）。
- SSE 尽早打开 + keep-alive 注释，降低「假挂起」误判。
- 客户端主动断开记为 disconnect，不算账号/上游失败。
- HTTP 入口 `recoverHandler` + keep-alive `safeCall`：panic 记日志返回 500，进程不退出。

#### 管理台与产品页
- 暗色 instrument 控制台（design tokens：`#0e1110` / `#e88a4a` / `#5ec4a8`）。
- `CODEBUDDY_PROXY_ADMIN_PASSWORD` 留空 = 本地免密；`/v1` API Key 仍强制可开。
- 一键国内 / 国际号池切换；区域端点以**账号**为准。
- 去掉 `?password=` query 鉴权；常量时间密钥比较；写操作 CSRF 同源校验。
- 产品页 / README / logo（Y 形信号路由）定稿；仓库更名为 `codebuddyapi-proxy`。

#### 开箱体验
- 二进制自动加载附近 `.env`；首次启动自动生成并持久化 API Key。
- 模型列表 60s TTL + singleflight 合并并发回源；`credits` 倍率完整透出。
- 采样参数 `temperature` / `top_p` / `max_tokens` 透传。

---

### 3. 内存与资源占用

- **运行时**：gz 实测 RSS 约 **5–7MB**（Go 单二进制，无 Node 运行时）。
- **磁盘抖动**：账号池高频路径不再每请求同步写盘，CPU/IO 尖峰消失。
- **仓库体积**：预编译包不再进 Git（历史 filter-repo 后 `.git` 约 **106MB → 12MB**）；发布只走 GitHub Releases。
- **分配**：SSE typed 解析、`Appendf` 组帧、惰性 prompt 估算，减少热路径临时对象与 GC 压力。
- **构建**：Makefile / release 增加 `-trimpath`，便于可复现构建。

---

### 4. 正确性与工程卫生（同版本一并收口）

- OAuth session 值快照 + 锁内鉴权（修 data race）。
- 运行时配置 `atomic.Pointer[config.Config]`（API Key / 号池切换并发安全）。
- 死代码清理、`strutil` 收敛、CI：`gofmt` + `go vet` + `go test -race`。
- `SECURITY.md`；核心包中文注释补全；适度关键路径测试（不追求堆用例）。

---

### 升级说明（v0.2.0 → v0.3.0）

1. 下载对应平台二进制覆盖即可；保留现有 `.env` 与 `proxy-accounts.json`。
2. systemd / 手工重启后生效；管理台密码若已清空则继续免密。
3. 下游若依赖 token 统计，请确认客户端读取流式末尾 usage chunk。
4. 强制刷模型：`GET /v1/models?fresh=1`。

### 资源清单

| 文件 | 平台 |
|------|------|
| `codebuddy-proxy-windows-amd64.exe` | Windows x64 |
| `codebuddy-proxy-linux-amd64` | Linux x64 |
| `codebuddy-proxy-darwin-arm64` | macOS Apple Silicon |
| `codebuddy-proxy-darwin-amd64` | macOS Intel |
| `SHA256SUMS.txt` / `.env.example` | 校验与配置模板 |

---

## 归档 · 2026-08-29 日间明细（已并入 v0.3.0）

以下条目已汇总进上方 v0.3.0，此处不再展开：审计 Batch A/B/C、OAuth/配置竞态、模型缓存与 credits、管理台免密与 CSRF、前端暗色 instrument、号池异步刷盘、token 透传、panic 兜底、热路径缓冲与 CJK 估算、仓库更名与 Pages 链接。

---

## 2026-08-28 · 首次启动 / 二进制体验

- 自动加载附近 `.env`（cwd / 可执行文件旁 / 上级目录）。
- 无 API Key 时首次启动自动生成并持久化到 `.env`，重启不丢。
- 默认 IDE 版本对齐 `2.117.2`。
- README / releases 补充上手说明。

---

## 历史里程碑（摘要）

- **Go 成为产品主线**：`protocol_direct` 网关；Pages：`main` + `/docs`。
- **主发布仓库**：`codebuddyapi-proxy`；`codebuddy-proxy` 作源码镜像。
- Node 版仅作本地历史备份，不再作为主发布路径。
