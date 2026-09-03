# 00 · Editorial 设计规格总纲（先读这份）

本目录是 CodeBuddy Proxy 前端改版的**唯一**设计源。  
**先读完这份再动手。** 旧版深色仪器风 / 北欧彩色 / 沙金+青+珊瑚 **全部作废**。

## 风格名

| 项 | 值 |
|---|---|
| style_name | Editorial |
| style_slug | editorial |
| 气质 | 高端杂志排版：衬线标题、无衬线正文、暖奶油底、软黑字、细线网格、大量留白 |

## 项目事实（决定实现边界，不决定审美）

| 项 | 事实 | 实现含义 |
|---|---|---|
| 用户 | 开发者，单人，本地机器 | 管理台仍是工具；审美走杂志，功能密度不能丢 |
| 网络 | 主要在国内 | **禁止任何 Web Font CDN** |
| 管理台部署 | 单文件 HTML（Go 反引号字符串） | 见下方硬约束 H1 |
| 产品页 | GitHub Pages 静态 HTML | 不增加 Go 二进制体积 |
| 产品 | Go 反代网关，OAuth + 账号池 | 文案继续说事，禁止营销黑话 |

## 冲突裁决（必须先读）

STYLEKIT 内部有两套互相打架的 token：

| 来源 | 说法 | 裁决 |
|---|---|---|
| Hard Prompt / Absolute Bans | **纯单色**，禁止红蓝绿等强调色 | **最高优先级，遵守** |
| Token Dictionary / Hero 模板 | `bg-[#e63946]` 红底 hero | **禁止使用** |
| Hard Prompt | 标题 `font-serif` weight 400，禁止 bold | 遵守 |
| Self-check | 禁止 Inter / Roboto / Geist / Fraunces / Plus Jakarta Sans | 遵守，系统字体栈 |
| Editorial Nav | `bg-[#F9F8F6]/90 backdrop-blur-sm` | 仅顶栏允许轻微 blur；**禁止**当作整页毛玻璃 |
| 旧 `design/`（深色 #0e1110、珊瑚 #e88a4a、薄荷 #5ec4a8） | 上一轮规格 | **作废，不得沿用** |
| worker 预览 `preview-admin-test.html` 北欧彩色 | 过期试验 | **作废，不得合入产品** |

**一句话：奶油底 + 软黑字 + 透明度层级。没有第三色。**

## 当前资产

| 资产 | 位置 | 本轮动作 |
|---|---|---|
| 产品页 | `docs/index.html` | 按 Editorial **重写视觉**；保留技术文案与下载信息 |
| Logo | `docs/logo.svg` | **单色重写**（见 `04-logo-and-motion.md`） |
| 管理台 | `internal/admin/page.go` 的 `PageHTML()` | 样式层重构为 Editorial；**id/JS/API 不动** |
| 落地页 | `internal/admin/page.go` 的 `LaunchPage()` | 跟随同一 token；极简 |
| 文档 | `README.md`、`docs/`、`CHANGELOG.md` | 同步产品页/Logo/管理台描述；**不改协议与 env 语义** |

## 硬约束（违反即返工）

### H1 · Go 反引号字符串边界

管理台 HTML 包在 `return \`...\`` 内。  
**即将写入的 CSS/JS 绝对不能出现反引号 \`。**  
JS 一律字符串拼接 `+`，禁止模板字符串。

### H2 · 禁止 Web Font CDN

删除全部 `fonts.googleapis.com` / `fonts.gstatic.com`。  
标题用系统衬线，正文用系统无衬线。禁止拉取 Inter / Fraunces / Outfit / Sora。

### H3 · DOM 契约不可破坏

见 `03-dom-contract.md`。  
**id 一律保持原样。** class 的**语义**（`.dot.bad`、`.badge.on`、`.seg button.active`、`.toast.show`）必须保留，只改 CSS。

### H4 · XSS 防护不可退化

`escapeHtml()` 必须保留且继续覆盖全部动态内容。  
`LaunchPage` 继续 `html.EscapeString`。

### H5 · 外部资源自包含

除 `docs/logo.svg` 被产品页/README 引用外，管理台不发任何外部请求。  
SVG 产品页可 `./logo.svg`；管理台 logo 必须 inline 或 data-URI。

### H6 · 不改后端

禁止改 `internal/server`、`internal/gateway`、`internal/accounts`、`internal/oauth`、`internal/provider`、路由、env、账号池协议。  
`page.go` 只动 HTML/CSS/JS 展示层。Go 函数签名、`PageHTML()` / `LaunchPage()` 出口保持。

### H7 · Editorial 禁令（审美最高优先级）

禁止：

- 彩色强调（红、蓝、绿、橙、青、珊瑚、沙金）
- `box-shadow` / `shadow-*`
- 粗边框（> 1px）
- 大圆角（`border-radius` ≥ 8px；本项目统一 **0**）
- 渐变、背景花纹、噪点 `feTurbulence`、装饰几何形
- 标题 `font-weight` ≥ 600
- 纯黑 `#000` / `#0a0a0a` 作主色
- 纯白 `#fff` / `#fafafa` 作背景
- 嵌套卡片
- 渐变文字、默认玻璃拟态
- Canvas 粒子 / 光柱 / glyph rain（与单色杂志冲突，旧 `04-product-ambient.md` 作废）

## 核心 token（三页同源）

```css
:root {
  --bg: #F9F8F6;
  --fg: #1C1C1C;
  --fg-80: rgba(28, 28, 28, 0.80);
  --fg-60: rgba(28, 28, 28, 0.60);
  --fg-40: rgba(28, 28, 28, 0.40);
  --fg-20: rgba(28, 28, 28, 0.20);
  --fg-10: rgba(28, 28, 28, 0.10);
  --invert-bg: #1C1C1C;
  --invert-fg: #F9F8F6;
  --display: Georgia, "Iowan Old Style", "Palatino Linotype", Palatino,
             "Songti SC", "Noto Serif SC", serif;
  --sans: system-ui, -apple-system, "Segoe UI",
          "PingFang SC", "Microsoft YaHei", sans-serif;
  --mono: ui-monospace, "SF Mono", Menlo, Consolas, monospace;
}
```

| 角色 | 用法 |
|---|---|
| 背景 | `#F9F8F6` |
| 主文字 | `#1C1C1C` |
| 次要 | `/60` |
| 辅助 / label | `/40` |
| 分割线 | `/10`，hover 可加深到 `/40` |
| 主按钮 | 反相：`background:#1C1C1C; color:#F9F8F6`（这不是第三色） |
| 次按钮 | 透明底 + 1px `/20` 底边或全边 |
| 状态（在线/失败） | **文字 + 透明度 / 实心 vs 空心**，不用绿/红 |

## 字体规则

- 标题：`var(--display)`，`font-weight: 400`，`letter-spacing` 更紧（tracking-tighter）
- Label：`var(--sans)`，`11–12px`，`letter-spacing: 0.2em`，uppercase，颜色 `/40`
- 正文：`var(--sans)`，`14–16px`，`line-height ≥ 1.625`，颜色 `/80`
- 装饰性副题：italic，`/60`
- 数字 / token / URL / 模型名：`var(--mono)`
- **禁止** Inter、Roboto、Geist、Fraunces、Plus Jakarta Sans、Outfit、Sora

## 动效

- 链接：underline 从右向左 `scaleX` 展开
- 标题：`group-hover` 转 italic，`transition 500ms`
- 图片（若有）：容器裁切 + 子元素 `scale(1.05)`，`1000ms`
- 分割线 hover：`/10` → `/40`
- 文字本身不位移
- `prefers-reduced-motion: reduce` 时关闭 transform/underline 动画
- **禁止** bounce、弹性曲线、整页骨架屏

## 执行顺序

```
README.md（本文件）     已读
        ↓
01-admin-console.md     管理台 + LaunchPage
        ↓
02-landing-page.md      产品页
        ↓
04-logo-and-motion.md   Logo SVG + 允许的动效
        ↓
03-dom-contract.md      全程对照，做完必跑
```

## 通用验收

```bash
# 1. 零外部字体
grep -rn "fonts.googleapis\|fonts.gstatic" docs/ internal/admin/page.go

# 2. 零噪点 / 零渐变文字
grep -rn "feTurbulence\|background-clip:\s*text" docs/index.html internal/admin/page.go

# 3. 零旧强调色
grep -nE "#e88a4a|#5ec4a8|#e86f3a|#e63946|#0e1110|#F4EEE4" \
  docs/index.html docs/logo.svg internal/admin/page.go

# 4. Go 未破坏
cd go-codebuddy && go build ./... && go vet ./... && go test ./...

# 5. 管理台无反引号模板字符串
grep -n '`' internal/admin/page.go   # 仅允许 Go raw string 的起止反引号
```

## 品味判据

靠约束，不靠「高端 / 精致 / 有质感」。  
截图去色后层级必须仍清晰。  
产品页一眼能看成杂志封面；管理台一眼能看成同一本杂志里的目录页，而不是另一套皮肤。
