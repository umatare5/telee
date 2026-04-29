# Repository Instructions

> [!IMPORTANT]
> Read [README.md](README.md) for project overview and [SECURITY.md](SECURITY.md) for security guidelines.

## Project Context

Mature production CLI for network automation. **Stability and backward compatibility are paramount** — CLI interface changes require explicit approval.

**Development Priorities**:

1. Stability over features
2. Minimal diffs
3. Named constants
4. Early returns

**Stack**:

- Go 1.25+ (see `go.mod`)
- Clean Architecture (`cmd/` → `cli/` → `internal/` → `pkg/`)
- Minimal dependencies (prefer standard library)
- Tooling: `golangci-lint`, `gotestsum`, `goreleaser` (see `Makefile`), `pre-commit` (see `.pre-commit-config.yaml`)

---

## Knowledge

Network devices have unique prompt patterns and privilege modes. Reference `internal/application/usecases/` and `internal/infrastructure/repositories/`:

- **Cisco IOS/IOS-XE**: Enable mode password required, prompt `>` → `#`
- **Cisco NX-OS**: Similar to IOS, different paging commands
- **Cisco ASA**: No paging support in user mode
- **Cisco AireOS**: No enable mode, WLC-specific prompts
- **Juniper SRX/SSG**: CLI vs config mode, prompts `>` / `#`
- **YAMAHA RT**: Japanese vendor, unique syntax
