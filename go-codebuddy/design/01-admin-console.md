# 01 · 管理台 Editorial 规格

**目标文件**：`go-codebuddy/internal/admin/page.go`  
- `PageHTML()`：管理台  
- `LaunchPage(title, message)`：OAuth 中间页  

**性质**：保留全部功能，替换设计语言。不是重写 JS 业务。

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

优先级：

```
可扫描性  >  杂志气场
信息仍在  >  留白仪式感
单色层级  >  状态彩点
```

产品页可以大留白；管理台 **同一套 token**，但 section 间距收紧，避免 1366×768 首屏只剩标题。

---

## §2 禁用清单（管理台）

| # | 禁止 | 处理 |
|---|---|---|
| E1 | 彩色 token（`--accent` `--teal` `--coral` `--warn` 作色值） | 删掉；状态改用 `/40` `/80` 实心/空心 |
| E2 | `feTurbulence`、多层 radial-gradient 光晕 | 删除 `body::before/::after` 装饰层 |
| E3 | 全页 `backdrop-filter` | 仅 `.topbar` 允许 `blur(8px)` + `background: rgba(249,248,246,.90)` |
| E4 | `border-radius` > 0 | 全部 `0` |
| E5 | `box-shadow` | 全部删除 |
| E6 | 按钮渐变 / inset 高光 | 主按钮纯反相填色 |
| E7 | Inter / Outfit 等 Web 字体 | 系统衬线 + 系统无衬线 |
| E8 | 胶囊 `border-radius: 999px` | 禁止；badge 改为 1px 方框 + uppercase label |
| E9 | Canvas / 粒子 | 管理台零 Canvas |

class 名 `.teal` `.danger` `.pill.good` **保留**（H3），只改 CSS：  
- `.primary` / `.teal`：反相填色 `#1C1C1C` / `#F9F8F6`  
- `.danger`：1px 边框 + 文字 `#1C1C1C`；hover 反相。删除用二次确认，不靠红色  
- `.dot.bad`：空心或 `/40` 方点；`.dot` 默认实心 `#1C1C1C`  
- `.badge.on` / `.off`：字重与边框对比，不用绿/灰彩

---

## §3 视觉规则

### 3.1 顶栏

```
position: sticky; top: 0; z-index: 50
background: rgba(249,248,246,.90)
backdrop-filter: blur(8px)          /* 仅此处 */
border-bottom: 1px solid rgba(28,28,28,.10)
height: 56–64px
左：衬线 logo 字  tracking 0.3em  uppercase  「CODEBUDDY」
右：健康文字 · 传输 · 站点    sans  xs  tracking 0.2em  uppercase  /60
```

不要胶囊容器包状态。用 `·` 分隔。

### 3.2 版式

```
最大宽度        1200px（管理台密度；不要用产品页那种 py-40）
左右 padding    24px / 48px
基础字号        14px
正文行高        1.625+
section 间距    32–48px（不是 96–160px）
卡片            1px border /10，padding 24px，radius 0
卡片 hover      border 加深到 /40，无阴影、不位移
```

区块标题：衬线，weight 400。  
Label：uppercase tracking 0.2em，`/40`。

### 3.3 表单

```
input/select: 无盒影、无 outline、无 ring
border: none; border-bottom: 1px solid rgba(28,28,28,.20)
font: 衬线 18–20px（短字段）；URL/Key 仍等宽
focus: border-bottom-color #1C1C1C
```

复制行（Base URL / API Key）允许等宽 + 底边，右侧「复制」为 hover-underline 文字按钮。

### 3.4 按钮

主操作（开始认证、拉取模型、生成 Key）：

```css
padding: 12px 24px;
font-size: 13px;
letter-spacing: .08em;
text-transform: uppercase;
background: #1C1C1C;
color: #F9F8F6;
border: 1px solid #1C1C1C;
border-radius: 0;
transition: color .2s, background .2s, border-color .2s;
```

次操作：透明底 + 1px `#1C1C1C`/20 边框，hover 底 `#1C1C1C`/0.02。  
删除：次操作样式，click 二次确认。

### 3.5 账号行

每条一行可扫：

```
标签  ·  站点  ·  登录态  ·  成功/失败  ·  Credits  ·  [查余额][刷新][禁用][删除]
```

数字等宽右对齐。操作默认 `/40`，hover 到 `/80`。  
禁止三颗实心彩按钮并排。

### 3.6 模型 chips

方框 1px `/10`，uppercase 或等宽模型 id。  
`.chip.btnish` 点击复制保留。hover 边框 `/40`，可 `font-style: italic`（500ms）。

### 3.7 动效

```
hover / focus     120ms
toast             300ms  底部偏右，仅 opacity
数据刷新          变化项背景闪 #1C1C1C/0.04，不做整页骨架
```

`:focus-visible`：2px `#1C1C1C` 描边 + 2px offset，不用光晕。

### 3.8 无障碍

- 状态不只靠颜色（本来就单色，必须有文字）
- 正文对比 ≥ 4.5:1（`#1C1C1C` on `#F9F8F6` 足够）
- `/40` 只用于 label，不用于关键数字
- 键盘可达

---

## §4 LaunchPage

OAuth 回调中间页：

1. 同一套 cream / ink token  
2. 无 Google Fonts、无圆角、无阴影、无 Canvas  
3. 居中：衬线标题 + sans 说明 + 一条返回管理台的 hover-underline 链接  
4. **保留** `html.EscapeString(title)` 与 `html.EscapeString(message)`  
5. 不要 py-40 仪式留白到看不清下一步

---

## §5 执行步骤

1. §0 文案  
2. 抽掉旧 `:root` 彩色 token 与 body 光晕  
3. 按 §3 重写 CSS；保留全部 id 与 JS 函数  
4. 重画账号行/chip/badge **仅 CSS + 必要 HTML 结构**，innerHTML 模板里的 class 语义不删  
5. LaunchPage 跟随  
6. `go build` / `go vet` / `go test`  
7. 对照 `03-dom-contract.md` 点功能  

---

## §6 验收

```bash
cd go-codebuddy
go build ./... && go vet ./... && go test ./...

grep -c "feTurbulence" internal/admin/page.go          # 0
grep -c "fonts.googleapis" internal/admin/page.go      # 0
grep -c "escapeHtml" internal/admin/page.go            # ≥ 10
grep -c "border-radius:999px" internal/admin/page.go   # 0
grep -cE "#e88a4a|#5ec4a8|#e63946" internal/admin/page.go  # 0
```

人工：

1. 1366×768 首屏能看到四项指标（登录 / Credits / 启用 / 成败）  
2. 3 秒内找到「刷新 Token」  
3. 断网打开，系统字体不塌  
4. OAuth / 号池 / 模型 / 复制 / 生成 Key 逐项可点  
5. 去色截图仍有层级（单色下必须靠字号/字重/线）  
6. 与产品页并排：同一本杂志，不是深浅两套皮  

---

## §7 失败模式

| 模式 | 表现 | 规避 |
|---|---|---|
| 换汤不换药 | 奶油底但仍留橙青光晕 | grep 旧色值 |
| 管理台做成产品页 | 指标被 py-40 推到三屏外 | §3.2 密度 |
| 用红表示删除 | 违反单色 | 二次确认 + 反相 hover |
| 改了 id | 按钮无反应 | H3 |
| 反引号 | 编译失败 | H1 |
| 漏 escapeHtml | XSS | H4 |
