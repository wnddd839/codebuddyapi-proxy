# Architecture Overview

## 目标

把 CodeBuddy 上游 `protocol_direct` 能力封装为 OpenAI 兼容网关：

```text
Client (OpenAI SDK / curl)
   │  /v1/chat/completions  (SSE or JSON)
   ▼
server ──► gateway ──► accounts pool + oauth refresh
                 │
                 ▼
              provider (protocol_direct)
                 │
                 ▼
         CodeBuddy upstream
         POST /v2/chat/completions
         GET  /v3/config
```

## 包职责

| 包 | 职责 | 不应包含 |
|----|------|----------|
| `config` | 环境变量解析 | 业务逻辑 |
| `accounts` | 账号 JSON 池、轮询选择、统计 | HTTP |
| `oauth` | plugin auth / refresh / JWT decode | 路由 |
| `provider` | 上游请求头、SSE 解析、事件累积 | 账号选择策略 |
| `models` | 模型发现与 public id 规范化 | 流式聊天 |
| `gateway` | 选号、刷新、失败重试、stats、OAuth session | HTML |
| `server` | 路由、鉴权、流式写回 | 上游协议细节 |
| `admin` | 管理台页面字符串 | 业务状态机 |
| `openai` | OpenAI chat/chunk 结构 | 上游映射 |
| `strutil` | `cmp.Or` 风格首个非空字符串 | — |
| `httputil` | JSON/SSE/Cookie/Origin | — |

## 设计原则

1. **单传输**：只维护 `protocol_direct`
2. **零第三方运行时依赖**（标准库优先）
3. **上下文可取消**：上游请求绑定 `context.Context`
4. **账号写盘原子化**：temp + rename
5. **现代 Go**：遵循 JetBrains `use-modern-go`

## 请求生命周期（chat）

1. API Key 校验
2. 解析 model → 上游 ID（兼容可选前缀 `codebuddy/`）
3. 从账号池选号（可按 site 过滤）
4. 必要时 refresh token
5. 组装 protocol_direct headers + body
6. 非流式：JSON；流式：SSE chunk + keep-alive
7. 成功/失败回写账号统计；鉴权失败可强制 refresh 后重试；站点非法可换号

## 与 Node legacy 的关系

`../legacy/nodejs` 仅归档。新功能不要回写 Node。
