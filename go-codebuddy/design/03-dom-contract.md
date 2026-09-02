# 03 · DOM 契约与回归清单（全程对照 · 改完必跑）

**这份是安全网。** 改 CSS 和 HTML 结构时，任何对 id/class 的重命名
都必须同步改 JS 与 Go —— 否则**功能静默失效**（不报错，就是点了没反应）。

> **强烈建议本次改版只动 CSS 与 HTML 结构，id 一律保持原样。**
> 重命名 id 的收益远小于出错风险。

---

## §1 硬约束回顾

| # | 约束 | 违反后果 |
|---|---|---|
| H1 | Go 反引号字符串内**禁止出现反引号 \`** | **编译直接失败** |
| H3 | id/class 是 JS↔HTML↔Go 三方契约 | 功能静默失效 |
| H4 | `escapeHtml()` 必须保留并覆盖全部动态内容 | **XSS 漏洞** |

### H1 详解（最容易踩）

```go
// page.go:9 开始，791 结束
return `<!doctype html>
...
</html>`
```

整个 HTML 是一个 Go **raw string literal**。因此：

```js
// ❌ 绝对禁止 —— 会提前终止 Go 字符串，编译失败
const html = `<div>${name}</div>`;

// ✅ 必须这样
const html = '<div>' + escapeHtml(name) + '</div>';
```

**当前状态：JS 里模板字符串数量 = 0。改版后必须保持 0。**

---

## §2 完整 id 清单（41 个 · 不可重命名）

### 状态与健康

| id | 用途 |
|---|---|
| `healthDot` | 健康指示灯（`.dot` / `.dot bad`） |
| `healthText` | 健康状态文字 |
| `pillTransport` | 传输方式 |
| `pillSite` | 当前站点 |
| `statusBox` | 状态区容器（默认 `display:none`） |
| `statusRaw` | 状态原始 JSON（`details.raw > pre`） |

### 核心指标（4 个 · 首屏必须可见）

| id | 用途 |
|---|---|
| `mLogin` | 登录态 |
| `mCredits` | Credits 余额 |
| `mEnabled` | 启用账号数 |
| `mSF` | 成功 / 失败计数 |

### OAuth

| id | 用途 |
|---|---|
| `site` | 站点选择（`<select>`） |
| `label` | 账号标签（`<input>`） |
| `btnStart` | 开始认证 |
| `btnPoll` | 检查登录 |
| `launchLink` | 打开登录页（`<a target="_blank">`） |
| `oauthBox` | OAuth 区容器 |
| `oauthMsg` | 会话状态消息 |
| `oauthRaw` | OAuth 原始 JSON |

### 账号池

| id | 用途 |
|---|---|
| `accounts` | 账号列表容器（JS 动态渲染） |
| `btnPoolDomestic` | 切换到国内站 |
| `btnPoolGlobal` | 切换到国际站 |
| `poolSiteSeg` | 号池分段控件容器 |
| `poolSiteHint` | 号池提示文字 |

### 模型

| id | 用途 |
|---|---|
| `modelChips` | 模型 chips 容器（动态渲染） |
| `modelsBox` | 模型区容器 |
| `modelsRaw` | 模型原始 JSON |
| `btnModels` | 拉取模型 |

### 客户端配置

| id | 用途 |
|---|---|
| `openAiBaseUrl` | Base URL |
| `openAiChatUrl` | Chat Completions URL |
| `openAiApiKey` | API Key |
| `openAiModel` | 推荐模型 |
| `copyBaseUrl` | 复制 Base URL |
| `copyChatUrl` | 复制 Chat URL |
| `copyApiKey` | 复制 API Key |
| `copyModel` | 复制模型名 |
| `btnGenerateKey` | 生成 API Key |
| `btnRefreshClient` | 刷新客户端配置 |
| `clientConfigHint` | 配置提示 |

### 其他

| id | 用途 |
|---|---|
| `btnRefresh` | 刷新状态 |
| `toast` | 全局 toast |
| `codebuddy` | OAuth 面板锚点 |

---

## §3 JS 函数清单（24 个 · 全部需保持可用）

| 函数 | 类别 | 备注 |
|---|---|---|
| `api` | 基础设施 | fetch 封装，`credentials:'same-origin'` |
| `escapeHtml` | **安全** | **H4，绝不可删** |
| `showToast` | UI | |
| `copyText` | UI | |
| `setHealth` | 渲染 | 用 `healthDot` + `healthText` |
| `fmtUptime` | 工具 | |
| `normalizeSite` / `siteLabel` | 工具 | 站点归一化 |
| `bareModelId` | 工具 | 去掉 `codebuddy/` 前缀 |
| `renderAccounts` | **渲染核心** | 动态生成账号卡片，大量 innerHTML |
| `renderModels` | **渲染核心** | 动态生成模型 chips |
| `paintStatus` | 渲染 | |
| `paintOAuth` | 渲染 | |
| `paintPoolSite` | 渲染 | |
| `paintClientConfig` | 渲染 | |
| `refreshStatus` | 数据 | |
| `refreshModels` | 数据 | |
| `refreshClientConfig` | 数据 | |
| `startOAuth` / `pollOAuth` | OAuth | |
| `switchPoolSite` | 操作 | |
| `onAccountAction` | 操作 | 账号行按钮事件委托 |
| `fetchAccountUsage` | 数据 | 查余额 |
| `generateApiKey` | 操作 | |

**事件绑定清单（13 个 onclick + 2 个 addEventListener）**：
`btnGenerateKey` `btnModels` `btnPoll` `btnPoolDomestic` `btnPoolGlobal`
`btnRefresh` `btnRefreshClient` `btnStart`
`copyApiKey` `copyBaseUrl` `copyChatUrl` `copyModel`
+ `addEventListener('change')` + `addEventListener('click')`（账号行委托）

---

## §4 CSS 类清单（53 个 · 可改样式，不可改语义）

**可自由调整视觉**：`.panel` `.panel-inner` `.metric` `.account` `.chip` `.badge`
`.pill` `.btn` `.ghost` `.primary` `.teal` `.danger` `.linkbtn` `.empty` `.err` `.meta`
`.mono` `.toast` `.field-grid` `.copyline` `.secret-hint` `.usage-line` `.chips`
`.actions` `.stack` `.section-head` `.eyebrow` `.lede` `.hero` `.topbar` `.brand`
`.mark` `.brand-copy` `.pillrow` `.shell` `.footer-note` `.oauth-status` `.seg`
`.seg-hint` `.metrics` `.account-top` `.account-title` `.raw` `.btnish`

**改样式时务必保留语义**：

| 类 | 语义 | 为什么不能乱改 |
|---|---|---|
| `.dot` / `.dot.warn` / `.dot.bad` | 健康状态 | JS 用 `dot.className` 切换 |
| `.badge.on` / `.badge.off` / `.badge.muted` / `.badge.site` | 账号状态 | `renderAccounts` 动态拼 |
| `.pill.good` / `.warn` / `.bad` | 余额等级 | `renderAccounts` 动态拼 |
| `.chip.btnish` | 可点击模型 | 绑定了 `data-copy` 委托 |
| `.seg button.active` | 号池选中态 | `paintPoolSite` 切换 |
| `.toast.show` / `.toast.err` | toast 显隐/错误 | `showToast` 切换 |

---

## §5 动态渲染的两处 XSS 边界（H4 · 重点）

以下两处用 `innerHTML` 渲染**服务端返回的数据**，所有插值必须过 `escapeHtml()`：

### 5.1 `renderAccounts`（约 `page.go:466`）

动态内容来源：账号 `label` / `userNickname` / `userName` / `userId` /
`site` / `authType` / `lastError` / `id` / 余额数值 / `officialUsageUrl`

**每一处插值都必须 escape**，尤其：
- `a.lastError`（上游错误消息，攻击者可部分控制）
- `usage.officialUsageUrl`（放进 `href`，注意 `javascript:` 协议风险）
- `a.id`（放进 `data-id` 属性）

### 5.2 `renderModels`（约 `page.go:535`）

动态内容来源：模型 `id` / `displayName` / `name` / `credits` / `description`

- 放进 `data-copy` 与 `title` 属性的都必须 escape
- `escapeHtml(tip || '点击复制')` 这类写法要保留

**验收**：改版后 `grep -c escapeHtml` 应 ≥ 10。低于此值说明有插值漏了。

---

## §6 回归测试清单（改完逐项点一遍）

自动化命令只能查语法和模式，**功能必须人工点**。

### 基础
- [ ] `go build` / `go vet` / `go test` 全过（验证 H1 未破坏）

### 状态与健康
- [ ] 页面加载后健康指示灯亮起（绿/红）
- [ ] 传输方式、站点 pill 显示正确
- [ ] 四项指标（登录态/Credits/启用账号/成功失败）有值
- [ ] 「刷新状态」按钮生效

### OAuth
- [ ] 选择站点 → 开始认证 → 生成登录链接
- [ ] 「打开登录页」新窗口跳转正确
- [ ] 「检查登录」轮询有响应
- [ ] 登录成功后账号自动进入账号池

### 账号池
- [ ] 账号列表正确渲染（含状态 badge、余额）
- [ ] 「查余额」返回 Credits
- [ ] 「刷新 Token」成功
- [ ] 「禁用 / 启用」切换生效
- [ ] 「删除」有二次确认且生效
- [ ] 号池切换（国内/国际）生效且分段控件高亮跟随

### 模型
- [ ] 「拉取模型」返回列表
- [ ] 点击 chip 复制模型名
- [ ] 原始 JSON（`details.raw`）可展开

### 客户端配置
- [ ] Base URL / Chat URL / API Key / 模型 显示正确
- [ ] 四个复制按钮均生效
- [ ] 「生成 API Key」成功并写入 .env
- [ ] 「刷新配置」生效

### 跨页
- [ ] 落地页（OAuth 回调）正常显示，返回链接可达

### 视觉
- [ ] 断网（无 Inter）时版式不塌
- [ ] 1366×768 下首屏四项指标完整可见
- [ ] 截图去色后层级清晰
- [ ] 产品页 → 管理台 → 落地页 观感连续

---

## §7 出问题时的定位顺序

```
1. 编译失败
   → 99% 是 H1：JS 里出现了反引号。搜 \` 定位。

2. 页面出来了但按钮点了没反应
   → H3：id 被改名。对照 §2 清单 diff 一遍。

3. 数据渲染成 [object Object] 或空白
   → §4 类名语义被改（如 .badge.on 被删），或 render 函数插值漏 escape 导致 HTML 破损。

4. 显示正常但功能全挂
   → 检查 addEventListener 的事件委托是否还在（账号行 / chips）。

5. 编译过、功能对，但就是看着不对
   → 回到 01 §2 禁用清单，跑 §6 的 grep 命令卡数字。
```
