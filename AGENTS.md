# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for project overview.

## Tech Stack

The list below covers the toolchain, the modules the binary links, and the release builder.

- Go 1.27+ (see [`go.mod`](go.mod))
- [`urfave/cli/v3`](https://github.com/urfave/cli) – flags, their `TELEE_*` sources and application lifecycle
- [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto) – `ssh` and `ssh/knownhosts`, the host-key path behind `pkg/ssh`
- [`goreleaser`](https://goreleaser.com/) – cross-platform release builds, configured by [`.goreleaser.yml`](.goreleaser.yml)

## Repository Structure

Read from [`cmd/main.go`](cmd/main.go) – each package is named for what it owns.

- [`cmd/`](cmd) – entry point, calling `cli.Start()` and nothing else
- [`cli/`](cli) – the one `cli.Command` and no subcommands, every flag declaration, and the `version` string ldflags sets
- [`internal/config/`](internal/config) – reads the flags into `Config`, and refuses a bad set before a socket opens
- [`internal/domain/`](internal/domain) – flag names, aliases, defaults and variables, the platform tokens, ports 22 and 23
- [`internal/application/`](internal/application) – one usecase per platform, forwarding to its repository
- [`internal/infrastructure/`](internal/infrastructure) – one repository per platform, owning the dialogue and the transport choice
- [`internal/framework/`](internal/framework) – routes `-x` to a usecase, and writes the transcript to stdout and a failure to stderr
- [`pkg/expect/`](pkg/expect) – the `BExp` and `BSnd` steps, and `Run`, which drives the login batch then the command batch
- [`pkg/telnet/`](pkg/telnet), [`pkg/ssh/`](pkg/ssh) – dial under `--timeout` and hand the connection to `expect.Run`
- [`pkg/errors/`](pkg/errors) – the sentinel errors the argument check returns
- [`docs/`](docs) – reference pages behind the README, indexed by [`docs/README.md`](docs/README.md)
- [`scripts/`](scripts) – helper scripts the pre-commit hooks run

## Setup and Commands

Run `make pre-commit-install` first.

- Read [`Makefile`](Makefile) for every make target and its requirements.
- Read [`CONTRIBUTING.md`](CONTRIBUTING.md) for the contribution rules.

## Code Style

Follow [Effective Go](https://go.dev/doc/effective_go) conventions and the software development principles DRY/YAGNI/SRP.

- Keep code simple and readable, avoiding clever tricks that obscure intent.
- Keep every change minimal, in code, tests, comments and documentation.
- Write simple comments that explain the reasoning behind the code, not just what it does.

## Testing

Follow [`CONTRIBUTING.md`](CONTRIBUTING.md).

- Run `make lint` and `make test-unit` before creating a commit.
- Take every fixture and sample identity from [`CONTRIBUTING.md`](CONTRIBUTING.md#fixture-identities), never a value read off a device.
- Place a test beside the code under test, in the `_test` package `testpackage` enforces.
- Tests cover `pkg/expect`, `pkg/telnet` and the pure functions of `pkg/ssh`, so raise the CI coverage threshold as tests land.

## Commits and PRs

Follow [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).

- Sign off every commit with `Signed-off-by:` (DCO).
- Open PRs against `main`. Create Draft PR as default.

## Domain Knowledge

Learn the constraints outside the CLI, because they decide what a batch has to absorb.

- Each platform is one usecase and one repository, named for the OS rather than the `-x` token, holding the whole session.
- `BExp` is a regular expression and `BSnd` is not, so the brackets AireOS and ScreenOS print need escaping.
- JunOS composes `username@hostname>`, so `--username` reaches the prompt match as well.
- IronWare alone needs `\r\n`, and every other platform takes `\n`.
- NX-OS alone prints a trailing space after the prompt character.
- ScreenOS alone prompts in lower case, printing `password:` where the rest print `Password:`.
- YAMAHA asks for no username, so `--username` reaches it over SSH only.
- JunOS has no paging command and takes `| no-more` on the command, so one already piped gains a second pipe.
- Every OS disables paging with a different command, and ASA needs a privileged session for it.
- See [`README.md`](README.md) for the platform matrix, and [`docs/troubleshooting.md`](docs/troubleshooting.md) for every refusal.
