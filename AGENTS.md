# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for the project overview, and [`docs/README.md`](docs/README.md) for the reference pages behind it.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`google/goexpect`](https://github.com/google/goexpect) v0.0.0-20210430 — the expect engine; a session is one `[]x.Batcher` of `BExp` and `BSnd` steps
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ — the single command, its thirteen flags, their aliases and their `TELEE_*` sources
- [`ziutek/telnet`](https://github.com/ziutek/telnet) v0.1 — the Telnet dialer `pkg/telnet` hands to goexpect's generic spawner
- [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto) v0.56+ — `ssh` and `ssh/knownhosts`, the host-key path behind `pkg/ssh`
- [`goreleaser`](https://goreleaser.com/) v2.18.0 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point; `main()` calls `cli.Start()` and carries nothing else
- `cli/` — The one `cli.Command` and no subcommands, every flag declaration, and the `version` string ldflags stamps
- `internal/config/` — Reads the flags into `Config` and runs `checkArguments`; a rejected set exits 1 before a socket opens
- `internal/domain/` — Flag names, aliases, defaults and env var names, the nine platform tokens, ports 22 and 23, the stderr hint
- `internal/application/` — One `Usecase` per platform, each forwarding to its repository so the router stays free of transport code
- `internal/infrastructure/` — One repository per platform, owning the whole wire dialogue and the Telnet-or-SSH choice
- `internal/framework/` — Routes `--exec-platform` to a usecase, writes device output to stdout and the error plus hint to stderr
- `pkg/telnet/`, `pkg/ssh/` — Dial, run one `ExpectBatch`, return the last match's output
- `pkg/errors/` — The sentinel validation errors `checkArguments` returns, quoted verbatim below

## Setup and Commands

Install required tools (one-time):

- `go install gotest.tools/gotestsum@latest`
- `golangci-lint` — See <https://golangci-lint.run/docs/welcome/install/>
- `gitleaks` — See <https://github.com/gitleaks/gitleaks#installing>
- `pre-commit` — See <https://pre-commit.com/#install>, then `make pre-commit-install` wires every hook in [`.pre-commit-config.yaml`](.pre-commit-config.yaml)

Make targets ([`Makefile`](Makefile)):

- `make build` — Build `tmp/telee` under `-trimpath` as the release build does, stamping `cli.version` from [`VERSION`](VERSION)
- `make lint` — `golangci-lint config verify` + `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make snapshot` — Build a `goreleaser` snapshot
- `make clean` — Remove build artifacts and `.bak*` files
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the hooks

The install passes `--allow-missing-config` because the hook path is the shared git common directory, so a hook installed from one worktree also fires in every other one and on `main`.

## Code Style

- [`.golangci.yml`](.golangci.yml) enables an explicit linter list over `default: none`, and `revive` caps a function at 80 statements
- Its `formatters` block runs `gci` and `gofumpt`, so `golangci-lint run --fix` rewrites formatting as well as logic
- A comment carries what the code cannot — a value another file must match, an order the device rejects
- One or two sentences, English, no emoji, and nothing a reader can derive from the code beside it

## Testing

- Run `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`), in the `_test` package `testpackage` enforces.
- Coverage is effectively zero: `cmd/main_test.go` asserts nothing, and both gates accept `0%`.
- The line above is temporary. Raise both gates as tests land, then delete it.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Committing on `main` is blocked by the `no-commit-to-main` hook, so branch first.
- Open PRs against `main`. CI runs the build and tests, lint, coverage, CodeQL, govulncheck, actionlint, markdownlint and the link check.
- A merged change to [`VERSION`](VERSION) tags and releases; nothing else triggers the release workflow.

## Domain Knowledge

Each platform is one usecase and one repository named for the OS rather than the `-x` token, and that repository holds the entire session. Nothing above it compensates for a prompt that fails to match.

What the devices do that a batch has to absorb:

- **`BExp` is a regular expression, `BSnd` is not.** AireOS and ScreenOS print brackets, so an unescaped one never matches.
- **JunOS composes `username@hostname>`**, so `--username` reaches the prompt match as well.
- **IronWare alone needs `\r\n`.** Every other platform takes `\n`.
- **NX-OS alone prints a trailing space** after the prompt character.
- **ScreenOS alone prompts in lower case**, printing `password:` where the rest print `Password:`.
- **YAMAHA asks for no username.** Its login opens on the password, so `--username` reaches it over SSH only.
- **JunOS has no paging command.** It takes `| no-more` on the command, so one already piped gains a second pipe.
- **Paging is disabled by a different command on every OS**, and on ASA that command needs a privileged session.

Read `internal/application/usecases/` and `internal/infrastructure/repositories/` before changing any of the above. [`README.md`](README.md) carries the platform matrix, [`docs/configuration.md`](docs/configuration.md) the flag surface, and [`docs/troubleshooting.md`](docs/troubleshooting.md) every refusal by its message.
