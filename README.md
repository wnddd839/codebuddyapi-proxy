# Proxy workspace (local integrator)

本目录是本地 **proxy 整合工作区**，不是直接推到 CodeBuddy GitHub 的仓库根。

```text
.
├── go-codebuddy/   ← CodeBuddy 反代（下次更新只推这个目录的内容到 GitHub）
├── go-cursor/      ← Cursor 反代（本地兄弟项目，gitignore，不随 CodeBuddy 发布）
├── .local/         ← 本地审查/草稿
└── scripts/        ← 推送辅助脚本
```

## 推送 CodeBuddy（只发 go-codebuddy）

```bash
# 预览将要发布的树（不推）
./scripts/push-codebuddy.sh --dry-run

# 推到 origin（codebuddyapi-proxy）的 main
./scripts/push-codebuddy.sh

# 同时推 mirror（codebuddy-proxy）
./scripts/push-codebuddy.sh --also-mirror
```

脚本用 `git subtree split --prefix=go-codebuddy`，远程仓库根目录会直接是 README / cmd / internal / docs，**不会**出现外层的 `go-cursor/`。

## 本地开发

```bash
cd go-codebuddy && make check
cd ../go-cursor && go test ./...
```
