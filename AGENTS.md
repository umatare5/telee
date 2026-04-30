# Repository Instructions

> [!IMPORTANT]
> Read [README.md](README.md) for project overview and [SECURITY.md](SECURITY.md) for security guidelines.

## Tech Stack

- Go 1.25+ (see [go.mod](go.mod))
- Clean Architecture (`cmd/` → `cli/` → `internal/` → `pkg/`)
- Minimal dependencies (prefer standard library)

## Repository Structure

- `cmd/` - Entry point (`main.go`)
- `cli/` - CLI flags and root command (urfave/cli/v3)
- `internal/application/` - Use cases per exec-platform
- `internal/domain/` - Domain entities and rules
- `internal/infrastructure/` - Repositories that drive telnet/ssh sessions
- `internal/framework/` - Cross-cutting framework code
- `internal/config/` - Runtime configuration
- `pkg/{telnet,ssh,errors}/` - Reusable transport and error packages

## Setup and Commands

Install required tools (one-time):

- `go install gotest.tools/gotestsum@latest`
- `golangci-lint` - See <https://golangci-lint.run/docs/welcome/install/local//>
- `goreleaser` release builds (see [.goreleaser.yml](.goreleaser.yml))
- `pre-commit install` wires `golangci-lint`, `gofmt`, `markdownlint-cli2`, `gitleaks` (see [.pre-commit-config.yaml](.pre-commit-config.yaml))

Make targets ([Makefile](Makefile)):

- `make build` - Build binary into `tmp/telee`
- `make lint` - `golangci-lint run` + `go mod tidy`
- `make test-unit` - run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` - generate HTML coverage report at `coverage/report.html`
- `make clean` - remove build artifacts

## Code Style

- `gofmt` and `golangci-lint` are enforced by the pre-commit hook (see [.pre-commit-config.yaml](.pre-commit-config.yaml)).
- Follow [.github/instructions/go.instructions.md](.github/instructions/go.instructions.md) for idiomatic Go.

## Testing Instructions

- Run `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`).
- Coverage threshold is enforced by [.github/workflows/go-test-coverage.yml](.github/workflows/go-test-coverage.yml).

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, tests, and CodeQL.

## Domain Knowledge

Network devices have unique prompt patterns and privilege modes.

Reference `internal/application/usecases/` and `internal/infrastructure/repositories/`:

- **Cisco IOS/IOS-XE**: Enable mode password required, prompt `>` → `#`
- **Cisco NX-OS**: Similar to IOS, different paging commands
- **Cisco ASA**: No paging support in user mode
- **Cisco AireOS**: No enable mode, WLC-specific prompts
- **Juniper SRX/SSG**: CLI vs config mode, prompts `>` / `#`
- **YAMAHA RT**: Japanese vendor, unique syntax

## References

- [.github/instructions/go.instructions.md](.github/instructions/go.instructions.md) - Idiomatic Go practices
- [.github/instructions/markdown.instructions.md](.github/instructions/markdown.instructions.md) - CommonMark (0.31.2)
- [.github/instructions/markdown-gfm.instructions.md](.github/instructions/markdown-gfm.instructions.md) - GitHub-flavored Markdown
- [.github/instructions/github-actions-ci-cd-best-practices.instructions.md](.github/instructions/github-actions-ci-cd-best-practices.instructions.md) - GitHub Actions CI/CD best practices
