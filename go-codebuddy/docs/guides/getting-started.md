# Getting Started

## 环境要求

- Go **1.26+**（`go.mod` 声明）
- 可选：已有 CodeBuddy OAuth 账号

## 安装与启动

```bash
cd go-codebuddy
cp .env.example .env
# 编辑 API Key / Admin Password / Site

go run ./cmd/codebuddy-proxy
# 或
make build && ./bin/codebuddy-proxy
```

## 最小配置

```env
CODEBUDDY_PROXY_HOST=127.0.0.1
CODEBUDDY_PROXY_PORT=32126
CODEBUDDY_PROXY_API_KEY=replace-with-a-long-random-key
CODEBUDDY_PROXY_ADMIN_PASSWORD=replace-with-a-long-admin-password
CODEBUDDY_SITE=domestic
CODEBUDDY_INTERNET_ENVIRONMENT=internal
```

## 首次登录

1. 打开 `http://127.0.0.1:32126/direct-admin/`
2. Basic Auth 使用 admin 密码
3. 选择 `domestic` 或 `global`，点击「开始 OAuth」
4. 浏览器完成授权后点「检查登录」
5. 账号写入 `CODEBUDDY_PROXY_ACCOUNTS_PATH`

## 验证

```bash
curl http://127.0.0.1:32126/health
curl http://127.0.0.1:32126/v1/models \
  -H "Authorization: Bearer $CODEBUDDY_PROXY_API_KEY"
```

## 常见问题

- **401 on /v1**：检查 `CODEBUDDY_PROXY_API_KEY`
- **no credentials**：先完成 OAuth
- **domestic chat endpoint**：国内站会走 `https://copilot.tencent.com/v2/chat/completions`
