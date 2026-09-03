# 04 · Logo 与动效（Editorial）

旧版 `04-product-ambient.md`（Canvas beams / glyph rain / 珊瑚薄荷粒子）**整份作废**。

本文件管两件事：

1. 重写 `docs/logo.svg`  
2. 产品页/管理台允许哪些动效  

---

## §1 Logo 重写

**文件**：`go-codebuddy/docs/logo.svg`  
**引用**：产品页 favicon / og:image、`README.md` 顶图、管理台 inline 一份同源图形。

### 概念（保持，只换颜色与线）

「多路进、/v1 出」：左侧折线汇入，右侧水平出线。  
这是协议翻译器，不是聊天气泡，不是字母 C 立体字。

### 新视觉

```svg
viewBox="0 0 256 256"
背景矩形  填 #F9F8F6    圆角 0（杂志，不要 iOS 图标圆角）
图形      描边 #1C1C1C  无渐变、无第二色
```

- 删除 `<linearGradient id="cbpIn">` 与 `#5EC4A8`  
- 描边宽约 14–16，`stroke-linecap="round"` 可保留（线端不是「圆角卡片」）  
- 可加极细十字或基线，对比度必须在 16×16 favicon 可辨  
- `role="img"` + `<title>CodeBuddy Proxy</title>` 保留  

**禁止**：彩色渐变、阴影滤镜、装饰圆点彩点、把 logo 做成红点杂志标。

管理台顶栏优先用 **衬线字标 CODEBUDDY**；若保留图形，必须 inline 同一路径，颜色改为当前 `--fg`，禁止再嵌珊瑚/薄荷 SVG。

---

## §2 允许的动效

| 效果 | 用在 | 参数 |
|---|---|---|
| hover-underline | 文字链接、次按钮 | `scaleX` 从右到左，200–400ms，`transform-origin: 100% 50%` |
| group-hover italic | 列表标题、模型名、账号标签 | 500ms |
| 分割线加深 | 网格/表行 | `/10` → `/40` |
| 反相填充 | 主按钮 hover 可保持填色，active `opacity .9` | 200ms |
| clip-path 揭示 | 仅产品页可选一张静图（无图则不做） | 不用库存 Unsplash |

### 禁止的动效

- Canvas / WebGL / 粒子 / 光柱 / 矩阵雨  
- bounce、elastic、整页 parallax  
- 图片 Ken Burns 超过 1.05 或快于 800ms  
- 文字位移（`translateY` 入场可在产品页 hero **一次**，管理台不要）  
- 管理台任何全屏氛围层  

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation: none !important;
    transition: none !important;
  }
}
```

---

## §3 体积与解耦

- Logo 与产品页动效 **不得** 写入 Go 依赖  
- 管理台零 Canvas，零外部 script  
- `docs/index.html` 增量不靠贴图；无大图资产  

---

## §4 验收

```bash
grep -nE "#e88a4a|#5EC4A8|#E86F3A|linearGradient" docs/logo.svg   # 空
grep -n "feTurbulence\\|canvas" docs/index.html                    # 空
```

人工：README 顶图、产品页 favicon、管理台字标三种尺寸都认得是同一个标。
