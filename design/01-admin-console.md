# 01 · 管理台设计规格（主任务 · 第一步执行）

**目标文件**：`go-codebuddy/internal/admin/page.go`
**当前**：853 行，其中 HTML/CSS/JS 全部包在 `return \`...\`` 内（第 9 行开，791 行收）
**性质**：保留全部功能，替换设计语言。**不是重写，是重构样式层。**

---

## §0 先修这两处（独立提交，5 分钟，与其他改动分开）

代码里残留了当初生成时的**指令文本**，正显示在所有用户的管理台首页：

```html
<!-- page.go:256 —— 用户可见的正文里写着"宣发计划"和"品牌主色策略" -->
<p class="lede">聚焦 OAuth 接入、账号池与请求健康。配色面向未来 GEO 站点宣发，可直接沿用品牌主色。</p>

<!-- page.go:407 —— 页脚备注同样是给设计师看的，不是给运维看的 -->
<p class="footer-note">Brand palette · sand / teal / coral · ready for GEO landing reuse</p>
```

**替换为**（面向终端用户，说人话）：

```html
<p class="lede">监控 OAuth 登录、账号池状态与请求健康度。</p>
<!-- 页脚那一整行直接删除，不留替代文本 -->
```

> 这两处是判断"这个界面有没有被认真对待"的最快信号。

---

## §1 它到底是什么（决定一切取舍）

一个**本地运行的开发者工具控制台**。开发者盯着它做这些事：

OAuth 登录 · 看账号池 · 查 Credits · 刷模型列表 · 复制 API Key/Base URL · 排查请求失败

**这是仪器，不是广告牌。** 三句话记住优先级：

```
信息密度  >  装饰
可扫描性  >  视觉冲击
稳定感    >  惊喜感
```

反例（现有实现全部命中）：一个本地工具，背景铺三层彩色光晕 + 噪点层 + 毛玻璃顶栏，
而真正要读的账号数据被挤在渐变卡片里。

---

## §2 禁用清单（违反任何一条即返工）

| # | 禁止 | 现有证据 |
|---|---|---|
| D1 | `feTurbulence` / SVG 噪点叠加层 | `body::before` 用了完整 feTurbulence data-URI |
| D2 | 三层以上背景径向渐变光晕 | body 背景叠了 3 层 radial + 1 层 linear |
| D3 | `backdrop-filter` 毛玻璃 | `.topbar` 用了 `blur(14px)` |
| D4 | `border-radius:999px` 用于超过 2 类元素 | **全文件 10 处**：按钮/badge/顶栏/分段器/chip/链接按钮… |
| D5 | 全局单一 `--ease` 套所有动效 | `--ease` 被 hover/focus/toast 共用，且统一 `.35s` |
| D6 | "沙金+青+珊瑚"暖中性模板色板 | `--sand #F4EEE4` / `--teal #0F7C74` / `--coral #E86F3A`，与产品语义无关 |
| D7 | 按钮用 `linear-gradient` + `inset` 高光组合 | `.primary` / `.teal` / `.danger` 三个都中招 |
| D8 | 按几何无衬线（Outfit）调版式却只用 `system-ui` 兜底 | 字体已降级，字距/行高/重心全是偏的 |

**D4 补充说明**：圆角必须有层级。当所有元素都是 999px，圆角就失去意义，只剩"软"。
规定：容器 `8px`，按钮 `6px`，**只有 badge 允许 999px**（且仅 1 类元素）。

**D5 补充说明**：真实设计里动效分层，绝不统一。规定三档：

```
微交互（hover / focus）  100–150ms  ease-out
组件级（展开 / 收起）    180–240ms
浮层（toast）            280–360ms
```

---

## §3 设计方向：工程仪器感

参考 **Linear / Vercel Dashboard / Tailwind UI 的工具面板**，**不是它们的营销页**。
关键词：**克制、对齐、可预测。**

### 3.1 颜色 token（深色，与产品页同源）

管理台目前是浅色沙金，产品页是深色编辑风 —— 用户点进来会以为进错了网站。统一为深色：

```css
:root{
  --bg:        #0e1110;   /* 与 docs/index.html 现有 --bg 完全一致 */
  --bg-elev:   #161b19;   /* 卡片 / 顶栏 */
  --bg-input:  #1c2220;   /* 输入框 */
  --ink:       #f2efe6;   /* 主文字 */
  --ink-dim:   #9aa39a;   /* 次要文字 */
  --line:      rgba(242,239,230,0.10);
  --line-soft: rgba(242,239,230,0.06);
  --accent:    #e88a4a;   /* 主操作 —— 沿用产品页 --accent */
  --teal:      #5ec4a8;   /* 成功 / 在线 */
  --warn:      #d4a24c;
  --danger:    #c9564a;
}
```

**彩色使用纪律**：
- 强调色只用于三处：**主操作按钮、状态指示、需要用户注意的数值**
- 其余一律走 `--ink` / `--ink-dim` / `--line` 三级灰阶
- **整个界面彩色像素占比 < 5%**（自测：截图去色，信息应基本不丢）

### 3.2 字体（先保证能用，再谈好看）

```css
--sans: "Inter", system-ui, -apple-system, "Segoe UI",
        "PingFang SC", "Microsoft YaHei", sans-serif;
--mono: ui-monospace, "SF Mono", "JetBrains Mono", Menlo, Consolas, monospace;
```

- **不引入任何 Web Font CDN**（见 H2）。Inter 缺失时退回 `system-ui` 完全可接受。
- 中文栈**必须**显式含 `PingFang SC`（mac）与 `Microsoft YaHei`（win）。
- 当前设计是照着 Outfit 的几何柔润感调的，而 Windows 的 `Segoe UI` 更宽更硬 ——
  **改版后必须重调字距与行长**，不能只换 font-family。

### 3.3 版式与密度

```
最大宽度        1200px
内容 padding    左右 24px
基础字号        14px
正文行高        1.6
次要信息        12px
卡片圆角        8px
卡片间距        16px
卡片内 padding  16px
```

**数字 / ID / Token / 模型名 / 端点 URL 一律等宽字体 + 右对齐** —— 便于纵向比对数值。

**状态指示**：不用色块 badge 时，用「纯文字 + 6px 圆点」，圆点**不带光晕描边**。

### 3.4 布局

```
顶栏  高 56px，不透明 --bg-elev，底部 1px --line
      左：品牌  右：在线状态 / 传输方式 / 站点
      用 "·" 分隔，不用胶囊容器（破 D4）

主区  两栏 1.2fr / 0.8fr（≥1024px），<1024px 单栏

账号列表 = 核心，给最大空间。每条压缩成一行内呈现：
      标签 · 站点 · 登录态 · 成功/失败 · Credits · [禁用][刷新][删除]

操作按钮默认低对比度，hover 才提亮 —— 避免三个 destructive 操作抢视觉。
删除类操作：hover 时边框转 --danger + 点击二次确认；不用实心红按钮常驻。
```

### 3.5 动效（只做有用的）

```
hover              背景色变化，120ms
focus-visible      2px --accent 描边 + 2px offset，不用 box-shadow 光晕
toast              底部偏右，位移 + 淡入，300ms
数据刷新           不做整页骨架屏动画；只让变化的那一项闪一次极短背景色
流式输出区         不做任何入场动画
```

### 3.6 无障碍

- 所有可点元素有 `:focus-visible` 且键盘可达
- 状态不只靠颜色传达（配文字或图标）
- 对比度：正文 ≥ 4.5:1，次要文字 ≥ 3:1

---

## §4 文案规则

**禁止**：赋能 / 一站式 / 极致 / 无缝 / 全新 / 打造 / 引领 / 驱动 / 生态 / 革新

**要求**：
- 直接说事。错误提示必须带**错误码 + 下一步动作**，不能只说"操作失败"。
- 现有文案里技术描述部分是好的（"OAuth 后直连上游，模型列表走 /v3/config"），保留这个风格。

---

## §5 执行步骤（建议顺序）

1. **§0 残片修复** → 单独一个 commit
2. **CSS 整体替换**：按 §3.1 token 重写 `:root` 与全部选择器，清 D1/D2/D3/D4/D6/D7
3. **字体栈**：换 §3.2，然后**重调字距/行高**（破 D8）
4. **动效分层**：拆掉全局 `--ease`，按 §3.5 三档改
5. **布局调整**：顶栏去毛玻璃、账号行压缩、操作按钮降对比度
6. **跑 §6 验收**

**全程保持 H1（不出现反引号）、H3（id 不变）、H4（escapeHtml 保留）。**

---

## §6 验收清单

```bash
cd go-codebuddy

# 硬约束
go build ./... && go vet ./... && go test ./...        # 必须全过

# D1 噪点层 = 0
grep -c "feTurbulence" internal/admin/page.go          # 期望 0

# D3 毛玻璃 = 0
grep -c "backdrop-filter" internal/admin/page.go       # 期望 0

# D4 胶囊 ≤ 2 处
grep -c "border-radius:999px" internal/admin/page.go   # 期望 ≤ 2

# D2 径向渐变 ≤ 1
grep -c "radial-gradient" internal/admin/page.go       # 期望 ≤ 1

# H2 零外部字体
grep -c "fonts.googleapis" internal/admin/page.go      # 期望 0

# H4 escapeHtml 仍在且被调用
grep -c "escapeHtml" internal/admin/page.go            # 期望 ≥ 10
```

**人工验收（必须做，命令查不出这些）**：

1. **断网字体降级**：注释掉所有 `@font-face` / 用无 Inter 环境打开，
   版式不塌、不重叠、可读性不受损（模拟国内网络）
2. **首屏不滚动可见**：登录态 / 启用账号数 / Credits / 失败数 四项全部可见
3. **3 秒法则**：一个陌生开发者能在 3 秒内找到「刷新 Token」按钮
4. **去色测试**：截图去色后信息层级依然清晰（证明不是靠颜色撑结构）
5. **功能回归**：OAuth 开始/轮询/回调、账号增删改、模型拉取、复制、号池切换 —— 逐项点一遍

---

## §7 常见失败模式（提前规避）

| 失败模式 | 表现 | 规避 |
|---|---|---|
| 换汤不换药 | 只改了配色，噪点/胶囊/渐变全留着 | 按 §6 命令逐条卡数字 |
| 破坏 DOM 契约 | 改名 id 后按钮点了没反应 | 见 H3，id 一律不动 |
| 反引号炸字符串 | Go 编译失败 | 见 H1，禁用模板字符串 |
| 降级后版式塌 | 无 Inter 时行高炸裂 | 见 §3.2，降级后重调 |
| XSS 退化 | 改模板时漏掉 escapeHtml | 见 H4 |
| 彩色过敏 | 深色底上又堆橙/青/红/黄 | 卡 §3.1「彩色 < 5%」纪律 |
