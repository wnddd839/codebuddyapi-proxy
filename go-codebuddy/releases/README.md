# Pre-built binaries

可直接下载运行，无需安装 Go。

| 平台 | 文件 |
|------|------|
| Windows x64 | `codebuddy-proxy-windows-amd64.exe` |
| Linux x64 | `codebuddy-proxy-linux-amd64` |
| macOS Apple Silicon | `codebuddy-proxy-darwin-arm64` |
| macOS Intel | `codebuddy-proxy-darwin-amd64` |

校验：`SHA256SUMS.txt`

## Windows 快速开始

1. 下载 `codebuddy-proxy-windows-amd64.exe`
2. 同目录复制 `.env.example` 为 `.env`，填写管理密码与 API Key
3. 双击或在终端运行：

```powershell
.\codebuddy-proxy-windows-amd64.exe
```

4. 打开 Admin：`http://127.0.0.1:32126/direct-admin/`
5. OpenAI 兼容 API：`http://127.0.0.1:32126/v1`

## Linux / macOS

```bash
chmod +x ./codebuddy-proxy-linux-amd64   # 或 darwin 对应文件
cp .env.example .env
./codebuddy-proxy-linux-amd64
```

## 环境变量（也可写入 `.env`）

| 变量 | 说明 |
|------|------|
| `CODEBUDDY_PROXY_ADMIN_PASSWORD` | 管理后台密码 |
| `CODEBUDDY_PROXY_API_KEY` | 客户端 API Key |
| `CODEBUDDY_PROXY_PORT` | 端口，默认 `32126` |
| `CODEBUDDY_SITE` | `domestic` 或 `international` |
| `CODEBUDDY_PROXY_ACCOUNTS_PATH` | 账号池 JSON 路径 |

完整说明见 [`docs/guides/getting-started.md`](docs/guides/getting-started.md)。
