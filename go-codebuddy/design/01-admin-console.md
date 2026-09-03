# 01 · 管理台 Editorial 规格

**目标文件**：`go-codebuddy/internal/admin/page.go`  
- `PageHTML()`：管理台  
- `LaunchPage(title, message)`：OAuth 中间页  

**性质**：保留全部功能，替换设计语言。严格执行 Editorial 杂志风与章节切页架构。

---

## §0 先修文案残片

若仍存在面向设计师的句子，改成给运维看的话：

```html
<p class="lede">监控 OAuth 登录、账号池状态与请求健康度。</p>
```

页脚品牌策略备注整行删除。禁止「宣发 / GEO / 品牌主色」出现在用户可见正文。

---

## §1 它是什么

本地开发者工具控制台，杂志排版。  
用户做：OAuth · 看号池 · Credits · 刷模型 · 复制 Base URL / API Key · 排错。

核心排版原则：

```
可扫描性  >  杂志气场
信息仍在  >  留白仪式感
单色层级  >  状态彩点
章节切页  >  单页堆叠
```

管理台采用 **Editorial 章节切页架构（Chapter Tabs）**，彻底避免将所有表单、账号卡片、模型列表全部单页垂直堆叠所导致的视觉臃肿。

---

## §2 禁用清单（管理台）

| # | 禁止 | 处理 |
|---|---|---|
| E1 | 彩色 token（`--accent` `--teal` `--coral` `--warn` 作色值） | 删掉；状态改用 `/40` `/80` 实心/空心 |
| E2 | `feTurbulence`、多层 radial-gradient 光晕 | 删除 `body::before/::after` 装饰层 |
| E3 | 全页 `backdrop-filter` | 仅 `.topbar` 允许 `blur(8px)` + `background: rgba(249,248,246,.90)` |
| E4 | `border-radius` > 0 | 全部 `0` |
| E5 | `box-shadow` | 全部删除（除切页当前 tab 底部内嵌 2px 指示线） |
| E6 | 按钮渐变 / inset 高光 | 主按钮纯反相填色 |
| E7 | Inter / Outfit / Fraunces 等 Web 字体 | 系统展示衬线 + 系统无衬线 |
| E8 | 胶囊 `border-radius: 999px` | 禁止；badge 改为 1px 方框 + uppercase label |
| E9 | Canvas / 粒子 | 管理台零 Canvas |

class 名 `.teal` `.danger` `.pill.good` **保留**（H3），只改 CSS：  
- `.primary` / `.teal`：反相填色 `#1C1C1C` / `#F9F8F6`  
- `.danger`：1px 边框 + 文字 `#1C1C1C`；hover 反相。删除用二次确认，不靠红色  
- `.dot.bad`：空心或 `/40` 方点；`.dot` 默认实心 `#1C1C1C`  
- `.badge.on` / `.off`：字重与边框对比，不用绿/灰彩

---

## §3 视觉与切页规则

### 3.1 顶栏

```
position: sticky; top: 0; z-index: 50
background: rgba(249,248,246,.90)
backdrop-filter: blur(8px)
border-bottom: 1px solid rgba(28,28,28,.10)
height: 56–64px
左：衬线 logo 字 tracking 0.2em uppercase 「CODEBUDDY」+ 传输协议
右：健康指示点与状态 · 传输方式 · 当前站点（sans xs tracking 0.15em uppercase /60）
```

### 3.2 章节切页导航（Chapter Navigation）

为消除单页垂直堆叠的臃肿感，管理台采用 4 大杂志章节切页：

```
[01 / 概览与监控]   [02 / 账号池与授权]   [03 / 客户端接入]   [04 / 模型与快照]
```

- **01 / 概览与监控 (`#tab-overview`)**：
  - 核心指标看板：登录态（`#mLogin`）、Credits（`#mCredits`）、启用数（`#mEnabled`）、成功/失败（`#mSF`）、总 Tokens（`#mTokens`）；
  - 刷新状态（`#btnRefresh`）与拉取模型（`#btnModels`）；
  - 429 智能熔断与 Reasoning 深度思考链双栏架构说明。
- **02 / 账号池与授权 (`#tab-pool`)**：
  - 账号池集群区域切换（`#poolSiteSeg`，`#btnPoolDomestic` / `#btnPoolGlobal`）；
  - 账号列表容器（`#accounts`）；
  - OAuth 授权卡片（`#codebuddy` 锚点，`#site`，`#label`，`#btnStart`，`#btnPoll`，`#launchLink`，`#oauthMsg`，`#oauthRaw`）。
- **03 / 客户端接入 (`#tab-client`)**：
  - OpenAI 兼容接入面板（`#client-config`，`#openAiBaseUrl`，`#openAiChatUrl`，`#openAiApiKey`，`#openAiModel`）；
  - 复制与生成 Key（`#copyBaseUrl`，`#copyChatUrl`，`#copyApiKey`，`#copyModel`，`#btnGenerateKey`，`#btnRefreshClient`）。
- **04 / 模型与快照 (`#tab-models`)**：
  - 可用模型芯片区（`#modelChips`，`#modelsRaw`）；
  - 实时快照 JSON 诊断（`#statusRaw`）。

**联动机制**：
- 支持 URL Hash 自动导航：
  - 访问 `/direct-admin/#codebuddy` 时（如 OAuth 成功后从 LaunchPage 点击返回）自动激活 **`02 / 账号池与授权`**；
  - 访问 `/direct-admin/#client-config` 时自动激活 **`03 / 客户端接入`**。
- **DOM 契约完全保留**：全部 41 个契约 ID 始终存在于 DOM 中，异步 API 轮询刷新零阻碍。

### 3.3 版式与间距

```
最大宽度        1200px
左右 padding    24px
基础字号        14px
正文行高        1.625
卡片            1px border /10，padding 24px，radius 0
卡片 hover      border 加深到 /40，无阴影、不位移
```

### 3.4 表单（底边线风格）

```css
input, select {
  width: 100%;
  padding: 8px 0;
  border: none;
  border-bottom: 1px solid rgba(28, 28, 28, 0.20);
  background: transparent;
  color: #1C1C1C;
  font-family: inherit;
  font-size: 14px;
  outline: none;
  transition: border-color 200ms ease;
}
input:focus, select:focus {
  border-bottom-color: #1C1C1C;
}
```

### 3.5 按钮

主操作（反相填色）：
```css
background: #1C1C1C;
color: #F9F8F6;
border: 1px solid #1C1C1C;
border-radius: 0;
font-size: 12px;
letter-spacing: .06em;
text-transform: uppercase;
```

次操作：透明底 + 1px 细边框，hover 加深。

### 3.6 动效纪律

```css
hover / focus    200ms
toast            300ms（底部偏右，仅 opacity 渐变，无移动）
切页淡入         250ms panelFade（从 opacity: 0 -> 1，微距 4px）
```

无障碍：`@media (prefers-reduced-motion: reduce)` 关掉所有动画。

---

## §4 LaunchPage

OAuth 回调中间页：
1. 纯单色 cream / ink token；
2. 无 Google Fonts、无圆角、无阴影、无 Canvas；
3. 居中：衬线标题 + sans 说明 + 返回管理台的反相链接（`/direct-admin/#codebuddy`）；
4. 保留 `html.EscapeString`。

---

## §5 验收

```bash
cd go-codebuddy
go build ./... && go vet ./... && go test ./...

grep -c "feTurbulence" internal/admin/page.go          # 0
grep -c "fonts.googleapis" internal/admin/page.go      # 0
grep -c "escapeHtml" internal/admin/page.go            # ≥ 10
grep -c "border-radius:999px" internal/admin/page.go   # 0
grep -cE "#e88a4a|#5ec4a8|#e63946" internal/admin/page.go  # 0
```
