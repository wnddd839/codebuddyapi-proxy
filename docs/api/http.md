# HTTP API

## Public

### `GET /health`

```json
{"ok":true,"provider":"codebuddy","transport":"protocol_direct"}
```

### `GET /v1/models`

Header: `Authorization: Bearer <API_KEY>`

返回 OpenAI list 形状。优先上游 `/v3/config`，失败回退配置中的 `auto`。

### `POST /v1/chat/completions`

Header: `Authorization: Bearer <API_KEY>`

Body（节选）：

```json
{
  "model": "glm-5.2",
  "stream": true,
  "messages": [{"role":"user","content":"hello"}]
}
```

- `stream=false` → `application/json`
- `stream=true` → `text/event-stream`，结束发送 `data: [DONE]`

兼容别名：`/models`、`/chat/completions`。

## Admin UI

- `GET /direct-admin/`
- Auth: Basic (`admin` / `CODEBUDDY_PROXY_ADMIN_PASSWORD`) 或 Bearer

## Admin API

| Method | Path | 说明 |
|--------|------|------|
| GET | `/direct-admin/api/status` | 运行状态 + 账号摘要 |
| GET | `/direct-admin/api/client-config` | 前端配置 |
| GET | `/direct-admin/api/codebuddy/status` | 账号池 |
| GET | `/direct-admin/api/codebuddy/accounts` | 账号列表 |
| DELETE | `/direct-admin/api/codebuddy/accounts/{id}` | 删除 |
| POST | `/direct-admin/api/codebuddy/accounts/{id}/enable` | 启用 |
| POST | `/direct-admin/api/codebuddy/accounts/{id}/disable` | 禁用 |
| POST | `/direct-admin/api/codebuddy/accounts/{id}/refresh-token` | 刷新 |
| GET | `/direct-admin/api/codebuddy/models` | 管理台模型列表 |
| POST | `/direct-admin/api/codebuddy/oauth/start` | 开始 OAuth |
| POST | `/direct-admin/api/codebuddy/oauth/poll` | 轮询 OAuth |
| GET | `/direct-admin/api/codebuddy/oauth/session` | 会话状态 |
| GET | `/direct-admin/codebuddy/oauth/launch` | 浏览器登录入口 |
| GET | `/direct-admin/codebuddy/oauth/callback` | 回调页 |

## CORS

`OPTIONS` 允许 `Authorization, Content-Type, X-API-Key`。
