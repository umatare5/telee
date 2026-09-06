# Contributing

Thank you for considering a contribution.

## Commands

The following `make` commands are available for development and testing:

| Command                     | Description                                              |
| :-------------------------- | :------------------------------------------------------- |
| `make help`                 | Display available targets and requirements               |
| `make build`                | Build the binary to `./tmp/telee`                        |
| `make lint`                 | Verify the lint config, run golangci-lint, tidy `go.mod` |
| `make test-unit`            | Run unit tests with coverage using gotestsum             |
| `make test-unit-coverage`   | Generate HTML coverage report                            |
| `make snapshot`             | Build a GoReleaser snapshot                              |
| `make clean`                | Remove build artifacts and backup files                  |
| `make pre-commit-install`   | Install the pre-commit hooks                             |
| `make pre-commit-test`      | Run every hook across the tree                           |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                              |

`make test-unit` clears the `TELEE_*` environment variables before running, because urfave reads them at flag-parse time and a developer's own shell would otherwise decide what the CLI tests see. Add any new variable the CLI reads to that list.

[`NOTICE`](NOTICE) reproduces the license of every module the binary links, so a change to the linked module set updates it.

Markdown style is enforced by the `markdownlint-cli2` hook that `make pre-commit-install` wires in, and again in CI. Links are checked in CI only, because that run reaches third-party hosts. Run `lychee .` to reproduce a link failure locally.

The hook path is the shared git common directory, so `make pre-commit-install` also arms every other worktree and the `main` checkout. It passes `--allow-missing-config` for that reason.

## Build

`make build` stamps `cli.version` from [`VERSION`](VERSION). That variable is initialised to `dev`, so a plain `go build ./cmd` still produces a working binary.

There is no target for the container image. GoReleaser builds it during a release and pushes it to `ghcr.io/umatare5/telee`, where a prerelease is excluded from the `latest`, `vX` and `vX.Y` tags.

## Release

To release a new version, follow these steps:

1. Rename `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md) to `## [vX.Y.Z]`, add that version's release link at the foot, and repoint the `[Unreleased]` compare link at the new tag.
2. Update the version in the [`VERSION`](VERSION) file.
3. Refresh the coverage badge — `make test-unit`, then `mkdir -p tmp && octocov badge coverage --config .octocov.yml > tmp/coverage.svg && mv tmp/coverage.svg docs/assets/coverage.svg`. The redirect goes through `tmp/` because the shell truncates its target before octocov runs. Nothing automates it — the reusable coverage workflow checks the threshold and writes no badge.
4. Submit a pull request with all three files.

Merging that pull request is the whole release. A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/telee/actions/workflows/go-release.yml), which tags the commit and publishes the release in the same run. The workflow has no manual trigger, so there is no step to perform by hand.

GoReleaser writes the release notes from the commit subjects. [`CHANGELOG.md`](CHANGELOG.md) is the hand-written counterpart and is not generated from them.

## Sample identities

Every transcript in this repository is device output with its identities replaced, and this section is the canon those replacements come from. It covers [`README.md`](README.md) and the reference pages under `docs/` equally.

Three kinds have a range reserved for exactly this:

| Kind         | Range                                     | Reserved by           |
| :----------- | :---------------------------------------- | :-------------------- |
| IPv4 address | `192.0.2.0/24`                            | RFC 5737              |
| MAC address  | `00:00:5e:00:53:00` – `00:00:5e:00:53:ff` | RFC 7042 §2.1.2       |
| Domain name  | `example.internal`                        | ICANN private-use TLD |

RFC 1918 space is excluded, because every address here sits on a `-H` a reader can paste and `192.168.0.0/16` may be live on that reader's own LAN. A transcript keeps the separator and case its device prints, which is why the same address reads `00-00-5E-00-53-01` under AlliedWare.

The rest have no standard behind them, so they are this repository's own:

| Kind         | Value           |
| :----------- | :-------------- |
| Switch       | `sw01` – `sw03` |
| Firewall     | `fw01`          |
| Controller   | `wlc01`         |
| Account name | `operator`      |

A password is never a literal. `<password>` and `<enable password>` stand in wherever a sample exports `TELEE_PASSWORD` or `TELEE_PRIVPASSWORD`, because a value that looks like it works is the one a reader copies. A vendor's own factory default is not an identity, so the AlliedWare sample keeps `-u manager`.

> [!IMPORTANT]
> Never paste a hostname, username, serial number or password from a capture into a sample transcript. Nothing in this repository redacts one, so a value pasted by hand reaches the tree unchanged.

## Pull requests

1. [Fork](https://github.com/umatare5/telee/fork) the repository
2. Create a feature branch
3. Commit your changes, using [Conventional Commits](https://www.conventionalcommits.org/) and signing off with `Signed-off-by:`
4. Add tests beside the code under test as `*_test.go`, and update the documentation alongside the code
5. Take every sample identity from [Sample identities](#sample-identities), never a value read off a device
6. Run `make lint` and `make test-unit`
7. Rebase your local changes against the `main` branch
8. Create a new Pull Request
