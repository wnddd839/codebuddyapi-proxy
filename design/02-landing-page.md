# 02 · 产品页 + 落地页设计规格（第二步执行）

**依赖**：先完成 `01-admin-console.md`，消费其 §3.1 颜色 token 与 §3.2 字体栈。

| 目标文件 | 行数 | 性质 |
|---|---|---|
| `docs/index.html` | 447 | 骨架健康，**只做减法** |
| `go-codebuddy/internal/admin/page.go:801-844`（LaunchPage） | 44 | 跟随统一 |

---

## §0 产品页现状判定：其实不差，别重做

我做过关键词扫描，AI 文案黑话（赋能/一站式/极致/无缝/全新/驱动/生态/打造/引领）
命中 **0 次**。现有文案是具体的、写给工程师的：

> "默认 protocol_direct。OAuth 后直连上游，模型列表走 /v3/config，不必再挂 CLI serve。"

这种诚实度很难得。**结构、配色策略、文案风格一律保留**，只做减法。

### 但也别被自己的判断骗了

产品页同样犯了两个错误，且与项目决策**自相矛盾**：

1. **还挂着 Google Fonts**（Fraunces + Sora）
   —— 你在管理台（`a5bb211`）已经移除过一次，产品页没跟上。
   本项目主要用户在国内，GitHub Pages 上加载 Fraunces/Sora 大概率失败或超时，
   直接退化成 Georgia/system-ui，还可能卡首屏。
2. **Fraunces + Sora 是当前 AI 生成产品页的头号固定搭配**
   （衬线展示体 + 几何无衬线）。一看到就有既视感。

---

## §1 必改（3 项）

### 1.1 移除 Google Fonts（H2）

删除全部 `fonts.googleapis.com` 与 `fonts.gstatic.com` 的 preconnect 与 stylesheet 引用。

替换为不依赖网络的字体栈：

```css
--display: "Iowan Old Style", "Palatino Linotype", Palatino,
           "Songti SC", Georgia, serif;     /* 有人文感的衬线，非 Fraunces */
--sans:    "Inter", system-ui, -apple-system, "PingFang SC",
           "Microsoft YaHei", sans-serif;   /* 非 Sora */
```

**关键**：降级后**必须重新检查行高与字距**。
Fraunces 的 x-height 和字重分布与 Iowan/Palatino 完全不同，
只换 font-family 不改排版，标题会散。

> 与 01 的 §3.2 保持完全一致（同一份 `--sans` 定义）。

### 1.2 删除 `.noise` 噪点层

`docs/index.html` 里的 `.noise` 用了完整 `feTurbulence`（与控制台 `body::before` 同源）。
这是 AI 生成页面的头号签名，且在这里纯装饰。

若确实需要质感，改用 **1px 网格线背景**（极低对比度）—— 真实技术文档站的做法。

### 1.3 收敛 hero 尺寸

```css
/* 现在 */
font-size: clamp(3.2rem, 10vw, 6.4rem);

/* 改为 */
font-size: clamp(2.6rem, 6.5vw, 4.2rem);
```

**理由**：现有尺寸导致首屏只剩标题。改后 hero-lead 与 CTA 在
**1366×768 笔记本上不滚动即可完整可见** —— 这是真实用户最常见的屏幕。

---

## §2 建议改（2 项）

- **`01/02/03` 编号栏**：改用图标，或去掉编号。编号栏是 AI 版式默认组件。
- **`.kicker`（"Only CodeBuddy"）**：保留，但降低字重 —— 它目前抢了 `h2`。

---

## §3 明确保留（不要动 —— 这些是对的）

- 深色底 + 双强调色（橙 `#e88a4a` + 青`#5ec4a8`）的配色策略
- 具体的技术文案风格（零营销黑话）
- 简洁页脚（项目 · 语言 · 许可 · 合规声明）
- 整体分区结构：hero / why / download / start / cta

---

## §4 与产品页统一（两页必须同源）

用户从产品页点进管理台时，不应感觉换了网站。强制执行：

| Token | 产品页 | 管理台 | 落地页 |
|---|---|---|---|
| `--bg` | `#0e1110` | `#0e1110` | `#0e1110` |
| `--accent` | `#e88a4a` | `#e88a4a` | `#e88a4a` |
| `--teal` | `#5ec4a8` | `#5ec4a8` | `#5ec4a8` |
| `--sans` | 01 §3.2 | 01 §3.2 | 01 §3.2 |

**产品页深色，管理台也深色。** 这是本轮改版最容易漏的一致性检查。

---

## §5 LaunchPage 落地页（`page.go:801-844`）

这是 OAuth 登录后的中间页，44 行，同样需要：

1. 移除 Google Fonts 引用（`Outfit`）
2. 应用 01 §3.1 的深色 token
3. 字体栈对齐 01 §3.2
4. 结构极简：居中卡片 + 状态文案 + （可选）一个返回管理台的链接
5. **保留 `html.EscapeString(title)` 与 `html.EscapeString(message)`** —— H4 安全边界

**不要**给它加噪点、渐变光晕或装饰图形。它只出现 2 秒，唯一任务是告诉用户"成功/失败了，接下来做什么"。

---

## §6 验收清单

```bash
# 零外部字体（三个文件）
grep -rn "fonts.googleapis\|fonts.gstatic" docs/index.html \
     go-codebuddy/internal/admin/page.go                      # 必须空

# 零噪点
grep -rn "feTurbulence" docs/index.html \
     go-codebuddy/internal/admin/page.go                      # 必须空

# Go 侧未破坏
cd go-codebuddy && go build ./... && go vet ./... && go test ./...
```

**人工验收**：

1. **断网降级**：禁用全部网络打开产品页，标题/正文排版不塌、不重叠
2. **首屏完整**：1366×768 下 hero 标题 + lead + 两个 CTA 全部可见，无需滚动
3. **跨页一致性**：产品页 → 管理台 → 落地页，背景色与字体观感连续（不像是三个网站）
4. **去色测试**：产品页截图去色后层级依然清晰
5. **对比检查**：与改版前并排看，确认"更像一个可信的开源项目"，而不是"更像一个模板"

---

## §7 常见失败模式

| 失败模式 | 表现 | 规避 |
|---|---|---|
| 自相矛盾 | 管理台无 Web Font，产品页还挂着 | §6 首条命令卡死 |
| 只换字体不调版式 | 标题散架、行高错乱 | §1.1 降级后重调，必做 |
| 误伤好文案 | 减法做过头，把技术细节也删了 | §3 明确保留清单 |
| 三页各做各的 | 三个页面三套色 | §4 表格逐行对 |
| LaunchPage 过度设计 | 给 2 秒的中间页加装饰 | §5 第 5 条 |
