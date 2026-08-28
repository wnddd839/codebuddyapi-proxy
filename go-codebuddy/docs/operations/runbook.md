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
