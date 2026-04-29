# Repository Instructions

> [!IMPORTANT]
> Read [README.md](README.md) for project overview and [SECURITY.md](SECURITY.md) for security guidelines.

## Tech Stack

- Go 1.25+ (see `go.mod`)
- Clean Architecture (`cmd/` → `cli/` → `internal/` → `pkg/`)
- Minimal dependencies (prefer standard library)

## Required Tools

- `golangci-lint`
- `gotestsum`
- `goreleaser` (see `Makefile`)
- `pre-commit` (see `.pre-commit-config.yaml`)

## Development Principles

- Stability over features
- Minimal diffs
- Named constants
- Early returns

## Project Knowledge

Network devices have unique prompt patterns and privilege modes. Reference `internal/application/usecases/` and `internal/infrastructure/repositories/`:

- **Cisco IOS/IOS-XE**: Enable mode password required, prompt `>` → `#`
- **Cisco NX-OS**: Similar to IOS, different paging commands
- **Cisco ASA**: No paging support in user mode
- **Cisco AireOS**: No enable mode, WLC-specific prompts
- **Juniper SRX/SSG**: CLI vs config mode, prompts `>` / `#`
- **YAMAHA RT**: Japanese vendor, unique syntax

## Additional Instructions

- [go.instructions.md](.github/instructions/go.instructions.md) - Instructions for writing Go code following idiomatic Go practices and community standards
- [markdown.instructions.md](.github/instructions/markdown.instructions.md) - Markdown formatting aligned to the CommonMark specification (0.31.2)
- [markdown-gfm.instructions.md](.github/instructions/markdown-gfm.instructions.md) - Markdown formatting for GitHub-flavored markdown (GFM) files
