# Documentation Index

Go CodeBuddy Proxy 文档目录。

## 全局

| 文档 | 说明 |
|------|------|
| [../README.md](../README.md) | 子项目总览与快速开始 |
| [guides/getting-started.md](guides/getting-started.md) | 安装、配置、首次 OAuth |
| [architecture/overview.md](architecture/overview.md) | 架构、包职责、请求链路 |
| [api/http.md](api/http.md) | HTTP API / Admin API |
| [operations/runbook.md](operations/runbook.md) | 运行、排障、资源占用 |
| [standards/coding-standards.md](standards/coding-standards.md) | 本仓库编码规范 |
| [standards/modern-go-guidelines.md](standards/modern-go-guidelines.md) | JetBrains 现代 Go 指南摘录（Go ≤ 1.26） |

## Vendor / Skills

| 路径 | 说明 |
|------|------|
| `.agents/skills/use-modern-go/` | JetBrains 官方 skill（含 CLI wrapper） |
| `.agents/skills/codebuddy-go/` | 项目级 skill |
| `standards/vendor/` | JetBrains FEATURES / LICENSE / guidelines.json 原文归档 |

## 阅读顺序建议

1. `guides/getting-started.md`
2. `architecture/overview.md`
3. `standards/coding-standards.md`
4. 需要改语法时加载 `use-modern-go` skill
