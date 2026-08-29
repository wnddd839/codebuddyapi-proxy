# 00 · 设计规格总纲（先读这份）

本目录是 CodeBuddy Proxy 前端改版的设计规格。**先读完这份再动手。**

## 项目事实（决定了所有设计决策）

| 项 | 事实 | 设计含义 |
|---|---|---|
| 用户 | 开发者，单人，本地机器 | 工具属性，不是营销页 |
| 网络 | 主要在国内 | **禁止任何 Web Font CDN** |
| 部署 | 单文件 HTML（Go 反引号字符串内） | 见下方「硬约束」 |
| 产品 | Go 反代网关，OAuth + 账号池 | 信息密度 > 装饰 |

## 当前状态

| 资产 | 位置 | 行数 | 判定 |
|---|---|---|---|
| 产品页 | `docs/index.html` | 447 | 骨架健康，需减法 |
| 管理台 | `go-codebuddy/internal/admin/page.go` | 853 | 需设计语言重构 |
| 落地页 | `internal/admin/page.go:801-844` | 44 | 需跟随统一 |

**关键发现：管理台现存 2 处生成指令残片**（详见 `01-admin-console.md` §0）。

## 硬约束（不可违反，违反即返工）

### H1 · Go 反引号字符串边界
管理台 HTML 全部包在 `return \`...\`` 内（`page.go:9` 开，`791` 收）。
**即将写入的 CSS/JS 绝对不能出现反引号 \`**。
现有 JS 恰好 0 处模板字符串 —— 保持这个状态，一律用字符串拼接 `+`。

### H2 · 禁止 Web Font CDN
删除全部 `fonts.googleapis.com` / `fonts.gstatic.com`。
产品页、管理台、落地页三处都要清。

### H3 · DOM 契约不可破坏
41 个 `id` + 53 个 `class` 是 JS 与 Go 的契约，见 `03-dom-contract.md`。
**重命名 id/class 必须同步改 JS 与 Go，否则功能静默失效。**
安全做法：只动 CSS 与 HTML 结构，保留 id 原样。

### H4 · XSS 防护不可退化
`escapeHtml()` 必须保留且继续覆盖全部动态内容（账号名 / 模型名 / 错误信息 / URL）。
这是安全边界，不是样式细节。

### H5 · 外部资源自包含
除 `.svg` 图标外不引入任何外部请求。SVG 一律 inline 或 data-URI。

## 执行顺序（有依赖，不可乱序）

```
01-admin-console.md   ← 第一步。确立 token，是后续的输入
        ↓
02-landing-page.md    ← 第二步。消费 01 的 token
        ↓
03-dom-contract.md    ← 全程对照，做完后跑验收
```

**理由**：产品页与管理台的 token 必须同源。反序会导致产品页定完色板后管理台还得返工。

## 通用验收（每份规格都适用）

```bash
# 1. 零外部字体
grep -rn "fonts.googleapis\|fonts.gstatic" docs/ go-codebuddy/   # 必须空

# 2. 零噪点层
grep -rn "feTurbulence" docs/ go-codebuddy/                       # 必须空

# 3. 零毛玻璃
grep -rn "backdrop-filter" go-codebuddy/internal/admin/page.go    # 必须空

# 4. Go 反引号字符串未被破坏
cd go-codebuddy && go build ./... && go vet ./...                 # 必须过

# 5. 单测仍绿（page.go 改动不得影响 Go 侧）
cd go-codebuddy && go test ./...
```

## 品味判据（贯穿全部规格）

> **靠约束传达品味，不靠形容词。**
> "高端 / 精致 / 有质感"这类词是 AI 味的主要来源，本目录一律不用。
> 需要什么效果，用可验证的规则描述（禁用清单 / 数值区间 / 验收命令）。

一个自测方法：**把成品截图去色，若层级依然清晰，说明结构成立；若变糊，说明是靠颜色撑的。**
