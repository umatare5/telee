# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for the project overview, and [`docs/README.md`](docs/README.md) for the reference pages behind it.

## Tech Stack

- Go 1.27 (see [`go.mod`](go.mod)); the four workflows that take a Go version pin 1.27.1
- [`google/goexpect`](https://github.com/google/goexpect) — the expect engine; a session is one `[]x.Batcher` of `BExp` matches and `BSnd` sends
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

- `make build` — Build `tmp/telee`, stamping `cli.version` from [`VERSION`](VERSION)
- `make lint` — `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run tests through `gotestsum` into `coverage/report.out`
- `make test-unit-coverage` — Render that profile to `coverage/report.html`
- `make clean` — Remove the binary, `tmp/dist`, `coverage/` and `.bak*` files, never `tmp/`, which holds the worktrees
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the hooks

The install passes `--allow-missing-config` because the hook path is the shared git common directory, so a hook installed from one worktree also fires in every other one and on `main`.

## Code Style

- [`.golangci.yml`](.golangci.yml) enables an explicit linter list over `default: none`, and `revive` caps a function at 80 statements
- It declares no `formatters` block, so `golangci-lint run --fix` rewrites logic and never formatting — the separate `gofmt` hook is what keeps the tree formatted, and it retires the day that block lands
- A comment carries what the code cannot: an order the device rejects, a value another file must match
- One or two sentences, English, no emoji, and never a restatement of the identifier beside it

## Testing

- Run `make test-unit` before committing.
- Place tests next to the code under test (`*_test.go`).
- Coverage is effectively zero. `cmd/main_test.go` is the only test file and it asserts nothing.
- The gates match that state: [`.github/workflows/go-test-coverage.yml`](.github/workflows/go-test-coverage.yml) passes `coverage_threshold: 0` and [`.octocov.yml`](.octocov.yml) accepts `0%`. Raise both as tests land.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Committing on `main` is blocked by the `no-commit-to-main` hook, so branch first.
- Open PRs against `main`. Every PR runs the build and tests, `golangci-lint` v2.13.2, coverage and CodeQL; actionlint runs only when one touches `.github/workflows/`.
- A merged change to [`VERSION`](VERSION) tags and releases; nothing else triggers the release workflow.

## Domain Knowledge

Each platform is a usecase and a repository under one package name — the OS name, not the `-x` token, so `srx` is `junos/` and `ssg` is `screenos/`. The repository holds the entire session, so editing one edits the wire dialogue — nothing above it compensates for a prompt that fails to match.

`-x` takes a token rather than the product name:

| `-x`      | Product                  | Transport   |
| :-------- | :----------------------- | :---------- |
| `aireos`  | Cisco AireOS             | Telnet, SSH |
| `allied`  | AlliedTelesis AlliedWare | Telnet only |
| `asa`     | Cisco ASA Software       | Telnet, SSH |
| `foundry` | Brocade IronWare         | Telnet only |
| `ios`     | Cisco IOS/IOS-XE         | Telnet, SSH |
| `nxos`    | Cisco NX-OS              | Telnet, SSH |
| `srx`     | Juniper JunOS            | SSH only    |
| `ssg`     | Juniper ScreenOS         | Telnet, SSH |
| `yamaha`  | YAMAHA RT                | Telnet, SSH |

`--secure-mode` on `allied` or `foundry` is refused with `secure-mode is not supported in this platform`, and omitting it on `srx` is refused with `non secure-mode is not supported in this platform`. `--default-privilege-mode` reaches only `asa`, `ios` and `nxos`; elsewhere it is refused with `default-privilege-mode is not supported in this platform`.

| `-x`      | Paging disabled by       | Privilege raised by |
| :-------- | :----------------------- | :------------------ |
| `aireos`  | `config paging disable`  | Not raised          |
| `allied`  | `terminal length 0`      | Not raised          |
| `asa`     | `terminal pager 0`       | `enable`            |
| `foundry` | `skip-page-display`      | `enable`            |
| `ios`     | `terminal length 0`      | `enable`            |
| `nxos`    | `terminal length 0`      | `enable`            |
| `srx`     | `\| no-more` appended    | Not raised          |
| `ssg`     | `set console page 0`     | Not raised          |
| `yamaha`  | `console lines infinity` | `administrator`     |

What bites when editing a batch:

- **`BExp` is a regular expression and `BSnd` is not.** AireOS matches `\(Cisco Controller\) >` and the ScreenOS redundant suffix is `\(M\)`, so an unescaped bracket never matches.
- **Seven prompts are built from `--hostname`.** AireOS is the exception with its fixed controller prompt, and JunOS composes `username@hostname>`, so `-H` must carry the configured hostname rather than an address.
- **IronWare is the only platform sending `\r\n`** — every other batch terminates a line with `\n` alone.
- **NX-OS is the only prompt carrying a trailing space**, so it matches `hostname>` and `hostname#` with one space appended.
- **ScreenOS is the only lower-case password prompt**, matching `password:` where the rest match `Password:`.
- **YAMAHA sends no username.** Its batch opens on `Password:`, so `--username` reaches the device over SSH only.
- **ASA and ScreenOS are the only platforms taking `--redundant-mode`**, which appends `/pri/act` and `\(M\)` to the prompt respectively.
- **JunOS sends no paging command.** It appends `| no-more` to the command itself, so a command already piped gains a second pipe.

Two guards decide whether a session starts at all:

- **ASA has no user-mode batch**, because `terminal pager 0` is privileged there, so a run with neither `-e` nor `-d` stops at `EnableMode must be set. Terminal length expansion in user-level is not supporting.`
- **Enable-mode is accepted but never acted on for `aireos`, `allied`, `srx` and `ssg`.** Those repositories build no privileged batch, so `checkArguments` prints `[INFO] enable-mode is ignored. It's not supported.` to stderr and the session continues unprivileged.

Read `internal/application/usecases/` and `internal/infrastructure/repositories/` before changing any of the above. Verified OS versions per platform sit in [`README.md`](README.md), and the full flag surface in [`docs/configuration.md`](docs/configuration.md).

> [!IMPORTANT]
> `--enable-mode` also demands that `--priv-password` differ from its default `enable`, or the run stops at `TELEE_PRIVPASSWORD must be set`. Pairing it with `--default-privilege-mode` stops at `enable-mode and default-priv-mode cannot use at once`.
