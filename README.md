# CodeBuddy Proxy · Go Subproject

独立封装的 **Go 版 CodeBuddy 纯协议反代网关**。

> 本目录是仓库的唯一实现。请在这里开发与运行。

## 快速开始

**直接下载：** [`releases/`](releases/) 内有 Windows / Linux / macOS 预编译包（见 [`releases/README.md`](releases/README.md)）。

**从源码：**

```bash
cp .env.example .env
go run ./cmd/codebuddy-proxy
```

```text
API:   http://127.0.0.1:32126/v1
Admin: http://127.0.0.1:32126/direct-admin/
Health: http://127.0.0.1:32126/health
```

## 项目结构

```text
go-codebuddy/
  cmd/codebuddy-proxy/     # 入口
  internal/                # 业务包
  docs/                    # 项目文档
  .agents/skills/          # 项目级 Agent Skills
    use-modern-go/         # JetBrains 官方 Modern Go skill
    codebuddy-go/          # 本仓库约定 skill
  Makefile
  README.md
```

## Agent Skills

本子项目内置两类 skill：

1. **`use-modern-go`**  
   来源：[JetBrains/go-modern-guidelines](https://github.com/JetBrains/go-modern-guidelines)（IDEA / GoLand / Junie 同源官方发布）  
   作用：按 `go.mod` 版本输出可使用的现代 Go 惯用法（`cmp.Or`、`slices.*`、`errors.AsType`、`new(expr)` 等）。

2. **`codebuddy-go`**  
   本仓库项目约定：包边界、禁止重新引入旧传输层、文档同步、安全日志等。

编辑 Go 代码前应先加载这两个 skill。

## 文档

完整文档索引见 [`docs/README.md`](docs/README.md)。

## 开发命令

```bash
make test
make build
make run
```

## License

BSD-3-Clause
