# Proxy workspace（本地整合目录）

外层只做整合；**推 GitHub 时只用 `go-codebuddy/`**。

```text
.
├── go-codebuddy/          CodeBuddy 反代（产品根：README / cmd / docs / CI）
├── go-cursor/             Cursor 反代（本地兄弟，不随 CodeBuddy 发布）
├── scripts/push-codebuddy.sh
├── .local/                私有部署笔记、审查稿、本地杂物
└── .tmp/                  临时克隆（如 9router 研究）
```

## 推送 CodeBuddy

```bash
./scripts/push-codebuddy.sh --dry-run     # 预览远程根目录
./scripts/push-codebuddy.sh               # 推 origin main（仅 go-codebuddy 内容）
./scripts/push-codebuddy.sh --also-mirror # 同时推 mirror
```

远程仓库根会是 `README.md`、`cmd/`、`internal/`、`docs/`……**不会**出现 `go-cursor/` 或本整合目录杂项。

## 本地开发

```bash
cd go-codebuddy && make check
cd ../go-cursor && go test ./...
```
