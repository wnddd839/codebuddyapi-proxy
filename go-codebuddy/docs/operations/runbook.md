# Operations Runbook

## 启动 / 停止

```bash
cd go-codebuddy
make build
./bin/codebuddy-proxy
```

信号：`SIGINT` / `SIGTERM` 触发 graceful shutdown（10s）。

## 健康检查

```bash
curl -fsS http://127.0.0.1:32126/health
```

## 日志

默认 JSON slog 到 stdout。关注字段：

- `codebuddy oauth credential refreshed`
- `retrying codebuddy request with next account`
- `retrying codebuddy request after oauth refresh`

## 排障

| 现象 | 排查 |
|------|------|
| listen 失败 | 端口占用；换 `CODEBUDDY_PROXY_PORT` |
| 401 API | API Key 不匹配 / `REQUIRE_API_KEY` |
| 无模型 | 未 OAuth；看 accounts path |
| 上游 401/403 | 点管理台 refresh-token 或重新登录 |
| 11140 / request illegal | site 与账号不匹配；检查 `CODEBUDDY_SITE` |
| 11128 / unapproved channel | 常见根因：客户端发了 `role=developer`（ZCode/OpenAI 新角色），上游会当成非法渠道。网关已映射为 `system`。也检查账号 `site=domestic` 时 endpoint 为 `copilot.tencent.com`（错误信息含 `[region=... endpoint=...]`） |
| 国内账号却打到海外 / 反过来 | 已按**账号 site** 选端点，不再看反代机器 IP，也不再被进程级 `CODEBUDDY_BASE_URL` 带跑偏 |
| 11101 / tool_choice unmarshal | 客户端传了对象型 `tool_choice`；网关会归一成字符串（`auto`/`none`/`required`） |
| 11102 / model service not found | 模型在列表里但上游无服务；换 `auto` 或其它可用模型 |
| Turn execution failed / reason=unknown | 看管理台 `stats.lastError`；多为上游 11128/11101 被客户端包装 |

## 资源占用参考

本地冒烟（空闲）：

- 二进制约 **10MB**
- Working Set 约 **10MB** 量级（显著低于 Node 常驻）

## 备份

定期备份 `CODEBUDDY_PROXY_ACCOUNTS_PATH` 指向的 JSON。


## OAuth「登录入口已失效」

根因是 session 对象被整体替换后，launch/callback 仍持有旧引用。

Go 版修复：
- `oauth` 使用指针并 `resetOAuthSessionLocked` 原地更新
- launch/callback 通过 `LiveOAuthSession()` 读取当前 session

验证：重新点「开始认证」，用**新生成**的 launch 链接，应 302 到 codebuddy.cn/login。
