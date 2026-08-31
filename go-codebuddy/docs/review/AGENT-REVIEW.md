# Agent Review — v0.3.7 Pool Cooldown Patch

**Status:** ✅ Blockers fixed (2026-08-31)  
**Version:** v0.3.7

## Blocker resolution

| # | Issue | Fix |
|---|-------|-----|
| 1 | All accounts cooling → `ErrNoAccounts` | `Select` fallback to min `CooldownUntil`; `Selection.BypassedCooldown` + debug log |
| 2 | `/model/info` 404 | Registered `GET /model/info` alias; `TestModelInfoAliasRoute` |

## Suggested fixes applied

| # | Issue | Fix |
|---|-------|-----|
| 3 | pre-commit silent gofmt | Changed to `gofmt -l` check only |
| 4 | Upsert migration | Documented in CHANGELOG v0.3.7 |

## Cooldown durations (heuristic, tunable)

| Error | Duration | Rationale |
|-------|----------|-----------|
| 11140/11128/11101/11102 | 5m | Policy errors — avoid hammering |
| 429 | 2m | Rate limit typical window |
| 502/503/504 | 30s | Transient upstream |

---

# Agent Review — Unreleased（冷却机制 + 规范落地）

**Branch/worktree:** local uncommitted  
**Date:** 2026-08-31  
**Scope:** Account pool cooldown, 11140 retry fix, Upsert site dedup, Go 1.26 hygiene, docs, Pi skill, pre-commit hook  
**Reviewer:** Any review agent (Bugbot / coding-worker-review-loop)

---

## Executive summary

| Area | Change | Risk |
|------|--------|------|
| Account pool | `cooldownUntil` + skip in `Select` | **Medium** — behavior change on failure paths |
| Gateway retry | 11140 removed from account rotation | **Low** — fixes false retries |
| Upsert | Dedup by `(userId + site)` not `userId` alone | **Medium** — existing users with collapsed accounts need re-OAuth |
| API surface | `MarkResult(..., cooldown)` signature | **Low** — internal only |
| JSON persistence | New optional field `cooldownUntil` on accounts file | **Low** — backward compatible (`omitempty`) |
| Docs | `/v1/model/info`, runbook WARN, architecture | **None** |
| Tooling | `.githooks/pre-commit` | **Low** — local dev only |
| Pi (out of repo) | `~/.pi/agent/skills/go-codebuddy-modern/` | **None** in git |

**Verdict request:** Confirm cooldown durations, 11140 no-retry, Upsert migration story, no JSON regression on account file load.

---

## Diff stat

```
11 files changed, 191 insertions(+), 18 deletions(-)
```

| File | Δ | Purpose |
|------|---|---------|
| `internal/accounts/pool.go` | +21 | Cooldown field, Select skip, MarkResult cooldown, Upsert site-aware |
| `internal/accounts/pool_test.go` | +64 | Cooldown + dual-site Upsert tests |
| `internal/gateway/service.go` | +35/-? | slices.Contains, failureCooldown, 11140 no retry |
| `internal/gateway/service_test.go` | +21 | failureCooldown + 11140 test |
| `internal/oauth/oauth.go` | +1 | `new(loggedIn)` |
| `docs/api/http.md` | +40 | GET /v1/model/info |
| `docs/operations/runbook.md` | +3 | WARN log level, cooldown note |
| `docs/architecture/overview.md` | ±4 | Failure handling table |
| `CHANGELOG.md` | +11 | Unreleased section |
| `.agents/*` | +8 | Pi system prompt + skill index |
| `.githooks/pre-commit` | new | gofmt + go test on staged Go |

Untracked (should be committed):

- `go-codebuddy/.agents/pi-system-prompt-go.md`
- `.githooks/pre-commit`

Not in git (manual install):

- `C:/Users/27297/.pi/agent/skills/go-codebuddy-modern/SKILL.md`

---

## Behavioral changes (review carefully)

### 1. Account cooldown

```go
// pool.go — Select skips when CooldownUntil > now
// MarkResult(false, err, cooldown) sets CooldownUntil = now + cooldown
// MarkResult(true, ...) clears CooldownUntil
```

| Error class | Cooldown | Still rotate account? |
|-------------|----------|------------------------|
| 11140 / 11128 / 11101 / 11102 | 5 min | **No** (11140/11128/11101/11102) |
| 429 / rate limit | 2 min | Yes (if another account exists) |
| 502 / 503 / 504 | 30 s | Yes |

**Edge case:** Single account in pool + cooldown → `ErrNoAccounts` until cooldown expires. **Intended** (stop hammering).

**Pinned account** (`AccountID` in SelectOptions): cooldown **not** checked — refresh/retry same account still works.

### 2. 11140 no longer triggers `retrying codebuddy request with next account`

Before: 11140 matched `reRetryNextAccount` → useless retry when only one global account.  
After: explicit `return false` in `shouldRetryNextAccount`; still applies 5 min cooldown via `failureCooldown`.

### 3. Upsert dedup `(userId + site)`

Before: same email OAuth on domestic then global **overwrote** domestic account.  
After: two separate pool entries.

**Migration:** Users who already lost domestic account must re-OAuth domestic. No automatic split of existing JSON.

---

## Go 1.26 / standards

| Item | Status |
|------|--------|
| `slices.Contains` for ExcludeIDs | ✅ Done |
| `new(loggedIn)` in oauth | ✅ Done |
| `slices`/`maps` elsewhere | ⏸ Not in scope |
| `errors.AsType[T]` | ⏸ No typed errors to unwrap yet |
| `omitzero` JSON migration | ⏸ **Deferred** — breaks `enabled: false` semantics |

---

## Tests run (must pass on review)

```bash
cd go-codebuddy
go test ./...
go build ./cmd/codebuddy-proxy
```

Expected: all packages `ok` (accounts, gateway, oauth, server, …).

New tests:

- `TestPoolSelectSkipsCooldown`
- `TestUpsertSameUserDifferentSites`
- `TestFailureCooldown`
- `TestShouldRetryNextAccount` — case `11140 no retry` → `false`

---

## Review checklist

### Correctness

- [ ] Cooldown cleared on successful `MarkResult(true)`
- [ ] `AccountID`-pinned Select bypasses cooldown (OAuth refresh path)
- [ ] All accounts in cooldown → `Select` returns `ErrNoAccounts` (not infinite loop)
- [ ] 11140 does not increment retry depth via account rotation
- [ ] Dual-site Upsert keeps distinct IDs and credentials

### Compatibility

- [ ] Old `proxy-accounts.json` without `cooldownUntil` loads fine
- [ ] No change to `/v1` JSON response shapes
- [ ] `MarkResult` signature change is internal-only (grep confirms)

### Security

- [ ] No token logging added
- [ ] pre-commit hook does not echo secrets

### Docs / DoD

- [ ] `docs/api/http.md` matches `server.go` route `GET /v1/model/info`
- [ ] runbook log level matches `main.go` default WARN
- [ ] architecture overview matches gateway retry logic

### Out of scope (do not block merge)

- omitzero migration
- oauth/admin test coverage
- pre-push use-modern-go CLI (pre-commit only runs go test)

---

## Suggested manual smoke (post-merge)

1. Start gateway, domestic + global OAuth both present.
2. Switch pool to global, chat `hy4-preview` — verify endpoint in error (if any) shows `codebuddy.ai`.
3. Force 429 (or mock) — confirm second request skips cooled account in admin `lastUsedAt` / rotation.
4. `GET /v1/model/info` with API key — returns LiteLLM-shaped `data[]`.

---

## Full patch

Generate locally:

```bash
git diff
git diff --stat
```

Or apply review on uncommitted tree at commit message:

```
fix(pool): cooldown isolation, 11140 no-retry, upsert by site

- Add cooldownUntil to account pool; skip cooled accounts in Select
- Stop rotating accounts on 11140; apply 5m cooldown
- Upsert dedup by userId+site (domestic/global coexist)
- slices.Contains, new(loggedIn), docs, pre-commit hook
```

---

## Review agent prompt (copy-paste)

```text
Review the unreleased patch in D:/codebuddy-proxy per go-codebuddy/docs/review/AGENT-REVIEW.md.

Focus:
1. Account cooldown correctness and single-account pool behavior
2. 11140 no-retry + cooldown interaction
3. Upsert (userId+site) migration risk
4. Test coverage for new behavior
5. Doc accuracy vs code

Run: cd go-codebuddy && go test ./...

Report: BLOCKERS / WARNINGS / NITs with file:line references.
Do not request omitzero migration in this patch.
```
