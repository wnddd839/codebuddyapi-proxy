# 04 · 产品页环境层（Ambient Layer）参考规格

**性质**：可选增强层，**只作用于** `docs/index.html`（GitHub Pages 静态页）。  
**不作用于**：管理台 `page.go`、LaunchPage、Go 二进制 —— **零运行时体积影响**。

**参考来源**：sysc-Go / sysc-greet 的终端动画**视觉语法**（beams、稀疏 glyph 雨、信号汇聚）。  
**禁止**：引入 sysc-Go、任何 npm/CDN 动画库、WebGL 重型引擎、位图字体 `.bit`。

**试验稿**：`docs/preview-particles.html`（对照用，合入正式页面前需过 §6 验收）。

---

## §0 为什么要写这份规格

产品页当前是「深色底 + 双强调色 + 静态 radial / 网格」——结构正确，但**环境偏静**，像固定容器。

sysc-Go 的启发不是「把 TUI 搬到 Web」，而是借三类**可验证**的氛围手段：

| sysc 效果 | Web 翻译 | 与产品的关系 |
|---|---|---|
| beams | 低对比斜向光柱，screen 混合 | 深度，不抢正文 |
| matrix / glyph rain | **稀疏** monospace 字符下落 | 开发者气质，密度远低于 TUI |
| （无直接同名） | **信号路由粒子**：双路珊瑚 → 单路薄荷 | 与 logo「多路进 /v1 出」同构 |

管理台保持 instrument 控制台风格：**不加全屏 Canvas**，避免与 `H3` 管理台「无 backdrop-filter 泛滥」冲突。

---

## §1 硬边界（与 Go 产品解耦）

| 规则 | 说明 |
|---|---|
| **A1** | 环境层代码只存在于 `docs/index.html`（及 preview 文件），**不**写入 `internal/admin/page.go` |
| **A2** | 不增加 Go 依赖、不增加 release 二进制体积 |
| **A3** | 单文件：Canvas + 内联 JS + 现有 CSS token，**无**外部 script/style |
| **A4** | `prefers-reduced-motion: reduce` 时**关闭** Canvas，回退纯 CSS 背景 |
| **A5** | 移动端（`max-width: 840px`）粒子数 ≤ 桌面 40%，或仅保留 grid + 静态 vignette |

> 用户关心的是 **Go 反代二进制占用低** —— 本层与 `codebuddy-proxy` 可执行文件无关；Pages 静态资源体增量目标 **< 25KB**（gzip 后 Canvas JS 控制在 ~8KB 内为佳）。

---

## §2 视觉 token（必须与 01/02 同源）

沿用 `02-landing-page.md` §4，**禁止**为「炫」引入第四主色。可选微量 accent：

```css
--bg: #0e1110;
--bg-elev: #161b19;          /* 或 rgba 玻璃层 */
--accent: #e88a4a;           /* 珊瑚：入站 / 光束暖色 */
--teal: #5ec4a8;             /* 薄荷：出站 / hub / 链接 */
--eldritch-muted: #7b68c4;   /* 可选：preview badge 级，opacity ≤ 0.35，不作正文色 */
--ink / --muted / --line      /* 与现网一致 */
```

**Eldritch 主题**（sysc `-theme eldritch`）仅作参考：深底 + 青绿/紫。**我们只借「深底上的微光」**，不照搬紫绿配色为主色。

---

## §3 允许的环境层（最多 4 层，可开关）

实现时每层独立布尔开关（便于 A/B 与降级），默认建议：

| 层 | 默认 | 实现要点 | 禁止 |
|---|---|---|---|
| **grid** | 开 | 72px 网格线，opacity ≤ 0.04，可选极慢 drift | feTurbulence / SVG 噪点 |
| **beams** | 开 | 6–8 束，宽 40–130px，globalCompositeOperation `screen`，opacity ≤ 0.06 | 全屏闪白、频闪 |
| **route** | 开 | 贝塞尔粒子双路汇入 hub，出站变 `--teal`；hub 径向 glow ≤ 90px | 粒子遮挡 hero 标题 |
| **glyph rain** | **关或极稀** | 字符集限定：`01/auto/v1/token/→` 等；≤ 48 滴；alpha 0.04–0.12 | 满屏 Matrix 绿、中文乱码 |

**内容层**：hero 使用 **半透明面板**（border + 低 blur 可选 ≤ 8px）浮于 Canvas 之上；section 正文仍用现有 rail / panel 结构，**不要**全站玻璃化。

---

## §4 原生实现约定（不搬库）

```text
结构：
  <canvas id="ambient" fixed, pointer-events:none, z-index:0>
  <main z-index:1> … 现有 DOM 不动 …

循环：
  requestAnimationFrame 单循环
  resize → 重算粒子数（按 w*h 上限 cap）
  每帧：clear → grid → beams → route → glyph → vignette

性能：
  - 不用 createImageBitmap / 离屏双缓冲（没必要）
  - 不用第三方粒子引擎
  - dpr = min(devicePixelRatio, 2)
  - 目标：桌面 idle < 3% 单核，60fps 可降级 30fps（移动端）
```

**与 sysc-Go 的差异**：sysc 在终端单元格上渲染；我们在 **CSS 像素 + Canvas 2D**。不移植 `.bit` 字体、不模拟 DOOM 火焰块。

---

## §5 与管理台 / 落地页的分工

| 页面 | 环境层 |
|---|---|
| `docs/index.html` | 可按本规格加 Canvas（可选） |
| Admin `page.go` | **仅** CSS 网格 / 1px line，**无** Canvas、**无** feTurbulence |
| LaunchPage | **无**装饰动画（§02 §5） |

跨页一致性靠 **token 同色**，不是靠同一套粒子脚本。

---

## §6 验收（合入 `index.html` 前必跑）

```bash
# 1. 仍零外部字体 / 零 feTurbulence
grep -rn "fonts.googleapis\|feTurbulence" docs/index.html   # 必须空

# 2. 无外部 script
grep -rn "<script src=" docs/index.html                    # 必须空

# 3. Go 不受影响
cd go-codebuddy && go build ./... && go test ./...

# 4. 文件体积（index 增量心里有数）
wc -c docs/index.html
```

**人工**：

1. 断网打开 Pages —— 动画与排版正常  
2. `prefers-reduced-motion` —— Canvas 隐藏，纯 CSS 可读  
3. 1366×768 —— hero 标题 + lead + CTA 无滚动可见  
4. 去色截图 —— 层级仍清晰（结构不依赖颜色）  
5. 与 `preview-particles.html` 并排 —— 正式版密度 **≤ preview**，更克制  

---

## §7 合入策略（给非前端维护者）

推荐 **两阶段**，避免一次改炸：

1. **Phase A（低风险）**：仅 hero 背后加 `grid + beams + route`，glyph rain 默认关  
2. **Phase B（可选）**：微调 hub 位置与品牌 aside（大屏 logo 半透明叠层）

不合入也可：现网 `index.html` 已达标；本规格是**增强路线**，不是 blocker。

**决策记录**（2026-08-29）：

- 参考 sysc-Go 氛围，**不**引入该库或任何终端依赖  
- 试验稿见 `docs/preview-particles.html`  
- Go 二进制体积与 RSS 目标不变；环境层仅静态页  

---

## §8 常见失败模式

| 模式 | 表现 | 规避 |
|---|---|---|
| TUI 廉价感 | 满屏绿字雨 | glyph 层默认关 / 极稀 |
| 性能债 | 笔记本风扇转 | §4 粒子 cap + reduced-motion |
| 容器更死 | 全站大圆角玻璃卡片 | 仅 hero 面板，section 保持 rail |
| 误伤二进制 | 把 Canvas 塞进 page.go | §1 A1 |
| 与管理台分裂 | 产品页很炫、admin 很素且色不同 | §2 token 表 + §5 分工 |
