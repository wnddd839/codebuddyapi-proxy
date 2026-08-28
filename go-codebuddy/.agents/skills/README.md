# Skills

This Go subproject vendors agent skills under `.agents/skills/`.

## Installed

### use-modern-go

- Source: JetBrains official repo `JetBrains/go-modern-guidelines`
- Purpose: force modern Go idioms matched to `go.mod`
- Wrapper scripts: `scripts/run-tool.ps1`, `scripts/run-tool.sh`

### codebuddy-go

- Source: this repository
- Purpose: package boundaries, docs sync, security, no legacy transports

## Update JetBrains skill

```bash
git clone --depth 1 https://github.com/JetBrains/go-modern-guidelines.git /tmp/go-modern-guidelines
cp /tmp/go-modern-guidelines/plugin/skills/use-modern-go/SKILL.md .agents/skills/use-modern-go/
cp /tmp/go-modern-guidelines/plugin/skills/use-modern-go/scripts/* .agents/skills/use-modern-go/scripts/
```
