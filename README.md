# Proxy workspace（本地整合目录）

外层只做整合；**推 GitHub 时只用 `go-codebuddy/`**。

```text
.
├── go-codebuddy/          CodeBuddy 反代（产品根：README / cmd / docs / CI）
├── scripts/push-codebuddy.sh
├── .local/                私有部署笔记、审查稿、本地杂物
└── .tmp/                  临时克隆（研究用）
```

## 推送 CodeBuddy

```bash
./scripts/push-codebuddy.sh --dry-run     # 预览远程根目录
./scripts/push-codebuddy.sh               # 推 origin main（仅 go-codebuddy 内容）
./scripts/push-codebuddy.sh --also-mirror # 同时推 mirror
```

远程仓库根会是 `README.md`、`cmd/`、`internal/`、`docs/`……**不会**出现本整合目录杂项。

## 本地开发

```bash
cd go-codebuddy && make check
```
