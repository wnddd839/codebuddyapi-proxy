<p align="center">
  <img src="docs/logo.svg" width="72" height="72" alt="CodeBuddy Proxy" />
</p>

<h1 align="center">CodeBuddy Proxy</h1>

<p align="center">
  <strong>把 CodeBuddy 变成 OpenAI 兼容的 <code>/v1</code> 渠道</strong><br/>
  协议直连 · OAuth 登录 · 账号池 · 管理台 · NewAPI 即插即用
</p>

<p align="center">
  <a href="https://wnddd839.github.io/proxy-codebuddy/"><img src="https://img.shields.io/badge/Product-Page-c45c26?style=flat-square" alt="Product Page" /></a>
  <a href="https://github.com/wnddd839/proxy-codebuddy/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-BSD--3--Clause-2f6b5a?style=flat-square" alt="License" /></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-%E2%89%A51.22-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" /></a>
  <img src="https://img.shields.io/badge/Transport-protocol__direct-3a4034?style=flat-square" alt="protocol_direct" />
</p>

<p align="center">
  <a href="https://wnddd839.github.io/proxy-codebuddy/">产品页</a> ·
  <a href="#快速开始">快速开始</a> ·
  <a href="#接入">接入 NewAPI</a> ·
  <a href="https://github.com/wnddd839/proxy-codebuddy">GitHub</a> ·
  <a href="CHANGELOG.md">更新日记</a>
</p>

---

## 这是什么

**CodeBuddy Proxy** 是一个专注 CodeBuddy 的自建反代（Go 原生实现）：

用你自己的 CodeBuddy 账号（OAuth），对外提供标准 OpenAI 兼容接口。  
管理台里加账号、看 Credits、刷模型；客户端 / NewAPI 只认 `/v1`。

> release传不上去，去产品页可以下载二进制打包。

仅限你拥有授权、且符合服务条款与合规要求的场景。  
**不要**把 token、管理密码、API Key 提交进仓库。

---

## 为什么用它

|  |  |
| :--- | :--- |
| **协议直连** | 默认 `protocol_direct`，OAuth 后直连上游，不依赖 `codebuddy --serve` |
| **真实余额** | 管理台拉官网 Credits，显示「剩余 / 总额」，不是假状态 |
| **模型列表** | 走协议 `/v3/config`，点一下刷新即可 |
| **账号池** | 多账号轮询；写盘防误清空；刷新失败只记错误 |
| **国内 / 国际** | 同一仓库；`.env` + 账号 `site` 切换，JWT 不可混用 |
| **OpenAI 形状** | `GET /v1/models` · `POST /v1/chat/completions` |
| **免安装运行** | 预编译二进制在 [`go-codebuddy/releases/`](go-codebuddy/releases/)，Windows 下载 exe 即用 |

---

## 快速开始

### 方式 A：直接下载（推荐，新手最省事）

1. 从 [产品页](https://wnddd839.github.io/proxy-codebuddy/) 或 [`go-codebuddy/releases/`](go-codebuddy/releases/) 下载对应平台二进制  
2. **直接运行**（同目录会自动读写 `.env`）  
3. 打开管理台登录 CodeBuddy，把页面上的 **Base URL + API Key** 填进 ZCode / NewAPI

Windows 示例：

```powershell
.\codebuddy-proxy-windows-amd64.exe
# 首次若没有 Key，日志会打印并写入同目录 .env
# 管理台：http://127.0.0.1:32126/direct-admin/
# API：   http://127.0.0.1:32126/v1
```

> **重要：** 管理台「生成 API Key」会**写入 `.env` 并立即生效**。旧 Key 会失效，客户端必须换成新 Key，否则会 401。

### 方式 B：从源码运行

```bash
git clone https://github.com/wnddd839/proxy-codebuddy.git
cd proxy-codebuddy/go-codebuddy
# 可选：cp .env.example .env
go run ./cmd/codebuddy-proxy
```

程序会自动加载当前目录 / 可执行文件旁的 `.env`（不会覆盖已有系统环境变量）。

可选手写配置：

```bash
CODEBUDDY_PROXY_API_KEY=你的长随机密钥
CODEBUDDY_PROXY_ADMIN_PASSWORD=你的管理台密码
CODEBUDDY_PROXY_REQUIRE_API_KEY=true
```

> 兼容旧版 Node 环境变量名（`CURSOR_DIRECT_*`），Go 版会一并读取。

| | |
| :--- | :--- |
| API | `http://127.0.0.1:32126/v1` |
| 管理台 | `http://127.0.0.1:32126/direct-admin/` |
| Health | `http://127.0.0.1:32126/health` |

打开管理台 → OAuth 登录 CodeBuddy → 复制 API Key 到客户端 → 开始调用。

管理台「账号池」可一键切换 **国内 / 国际** 号池：切换后只请求该区域账号，并写入 `.env`（`CODEBUDDY_SITE`）。

---

## 接入

**任意 OpenAI 兼容客户端 / NewAPI / Sub2API**

```text
Base URL   http://<host>:32126/v1
API Key    与 .env 中 CODEBUDDY_PROXY_API_KEY 相同
Model      codebuddy/auto  或 GET /v1/models 返回的 id
```

下游可直接拉模型列表（标准 OpenAI Models API）：

```bash
# 列表
curl http://127.0.0.1:32126/v1/models \
  -H "Authorization: Bearer $CODEBUDDY_PROXY_API_KEY"

# 单个模型
curl http://127.0.0.1:32126/v1/models/codebuddy%2Fauto \
  -H "Authorization: Bearer $CODEBUDDY_PROXY_API_KEY"
```

返回形如：

```json
{
  "object": "list",
  "data": [
    {
      "id": "codebuddy/auto",
      "object": "model",
      "created": 1700000000,
      "owned_by": "codebuddy"
    }
  ]
}
```

有 OAuth 账号时，列表会尽量来自协议 `/v3/config`；国际站 `/v3/config` 常空时会合并 CLI 模型目录（gpt / gemini 等）。否则至少返回 `codebuddy/auto`（上游映射为 `default-model`）。

```bash
curl http://127.0.0.1:32126/v1/chat/completions \
  -H "Authorization: Bearer $CODEBUDDY_PROXY_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"codebuddy/auto","messages":[{"role":"user","content":"你好"}]}'
```

---

## 国内 / 国际

同一仓库，靠 `.env` + 账号 `site` 切换；**JWT 不可混用**。

国内（默认）：

```text
CODEBUDDY_SITE=domestic
CODEBUDDY_INTERNET_ENVIRONMENT=internal
CODEBUDDY_BASE_URL=https://www.codebuddy.cn
```

国际：

```text
CODEBUDDY_SITE=global
CODEBUDDY_INTERNET_ENVIRONMENT=public
CODEBUDDY_BASE_URL=https://www.codebuddy.ai
```

注意：

- 国际「活动赠送包」常能看积分但 chat 报 `11140` → 换 Pro 试用/付费号
- Sub2API 只发 `"hi"` 探活：网关会自动补 system，避免 `11101`
- `11101` 不当作换号；`11140` 多半是套餐权限

---

## 文档

| 资源 | 链接 |
| :--- | :--- |
| 产品页 | [wnddd839.github.io/proxy-codebuddy](https://wnddd839.github.io/proxy-codebuddy/) |
| 文档索引 | [`go-codebuddy/docs/README.md`](go-codebuddy/docs/README.md) |
| 预编译包 | [`go-codebuddy/releases/README.md`](go-codebuddy/releases/README.md) |
| 编码标准 | [`go-codebuddy/docs/standards/coding-standards.md`](go-codebuddy/docs/standards/coding-standards.md) |

常用命令（仓库根目录）：

```bash
make test    # go test ./...
make build   # 产出 go-codebuddy/bin/
make run
```

---

## License

[BSD-3-Clause](LICENSE)

灵感与参考：[CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) · [Kiro-Go](https://github.com/Quorinex/Kiro-Go) · NewAPI 生态
