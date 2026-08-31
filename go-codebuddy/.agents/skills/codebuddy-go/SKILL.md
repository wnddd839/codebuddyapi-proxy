---
name: codebuddy-go
description: Project conventions for the Go CodeBuddy protocol proxy in go-codebuddy. Use when editing Go code, docs, routes, OAuth, account pool, or provider streaming in this subproject.
---

# CodeBuddy Go Subproject Skill

This repository's supported implementation lives in `go-codebuddy/`.

## Required companion skill

Before writing or refactoring Go in this project, also apply the official JetBrains skill:

- `.agents/skills/use-modern-go` (from [JetBrains/go-modern-guidelines](https://github.com/JetBrains/go-modern-guidelines))

**Pi agents:** load global skill `go-codebuddy-modern` (`~/.pi/agent/skills/`) — it wraps both skills and the mandatory CLI workflow. System prompt: `.agents/pi-system-prompt-go.md`.

On Windows:

```powershell
.\.agents\skills\use-modern-go\scripts\run-tool.ps1 list --file-path path\to\file.go
```

On Unix:

```bash
sh .agents/skills/use-modern-go/scripts/run-tool.sh list --file-path path/to/file.go
```

Target Go version is declared in `go.mod` (currently **1.26+**). Prefer idioms from JetBrains guidelines up to that version.

## Architecture boundaries

| Package | Responsibility |
|---------|----------------|
| `cmd/codebuddy-proxy` | process entry / signals |
| `internal/config` | env loading only |
| `internal/accounts` | account pool persistence |
| `internal/oauth` | plugin OAuth + refresh |
| `internal/provider` | protocol_direct upstream + SSE |
| `internal/models` | `/v3/config` discovery |
| `internal/gateway` | orchestration (select/refresh/retry/stats) |
| `internal/server` | HTTP routes |
| `internal/admin` | admin HTML |
| `internal/openai` | OpenAI response shaping |
| `internal/strutil` | shared string helpers (`cmp.Or`) |
| `internal/httputil` | JSON/SSE/auth helpers |

Do **not** reintroduce `cli_daemon` / ACP / cloud multi-transport as default paths.

## Coding rules for this repo

1. Keep zero third-party runtime dependencies unless absolutely required.
2. Use Go 1.22+ method-aware `ServeMux` patterns for new routes.
3. Prefer `strutil.First` / `cmp.Or`, `strings.CutPrefix`, `min`/`max`, `slices.*`, `maps.*`.
4. Account pool writes must stay atomic (`*.tmp` + rename).
5. Never log full bearer/refresh tokens; only previews/hashes.
6. Streaming handlers must support client cancel via request context.
7. Update docs under `docs/` when changing routes, env vars, or package boundaries.
8. Add/adjust tests for account pool, config defaults, and provider parsing when touching those areas.

## Docs map

- Global intro: `README.md`
- Doc index: `docs/README.md`
- Architecture: `docs/architecture/overview.md`
- API: `docs/api/http.md`
- Ops: `docs/operations/runbook.md`
- Coding standards: `docs/standards/coding-standards.md`
- JetBrains vendor reference: `docs/standards/modern-go-guidelines.md`

## Definition of done

- `go test ./...` passes
- `go build ./cmd/codebuddy-proxy` succeeds
- New public behavior documented in `docs/`
- Modern Go guidelines considered via `use-modern-go`
