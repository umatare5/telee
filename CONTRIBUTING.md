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
| `make test-unit-coverage`   | Generate the HTML coverage report                        |
| `make snapshot`             | Build a GoReleaser snapshot                              |
| `make clean`                | Remove build artifacts and backup files                  |
| `make pre-commit-install`   | Install the pre-commit hooks                             |
| `make pre-commit-test`      | Run every hook across the tree                           |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                              |

`make test-unit` clears the `TELEE_*` environment variables before running, because urfave reads them at flag-parse time and a developer's own shell would otherwise decide what the CLI tests see. Add any new variable the CLI reads to that list.

`make clean` names the binary, `tmp/dist` and the generated HTML report one by one rather than removing `./tmp` or `./coverage`. Removing either directory would take out the worktrees `git wt` puts under `./tmp` and the tracked `coverage/report.out`.

`make pre-commit-install` passes `--allow-missing-config`, because the hook path is the shared git common directory. A hook installed from one worktree fires in every other one and on the `main` checkout, where `.pre-commit-config.yaml` may not exist.

Five hooks run on a commit — `no-commit-to-main`, `golangci-lint-full`, `actionlint`, `gitleaks-system` and `markdownlint-cli2`. The `golangci-lint` `rev` in [`.pre-commit-config.yaml`](.pre-commit-config.yaml) has to match `golangci_lint_version` in [`.github/workflows/go-test-fmt.yml`](.github/workflows/go-test-fmt.yml). Renovate updates the two through different datasources, so nothing forces them onto the same pull request.

No separate `gofmt` hook sits beside them, because `.golangci.yml` enables `gci` and `gofumpt` under `formatters`. `golangci-lint run` reports a badly formatted file the way it reports any other finding, so the `golangci-lint-full` hook already fails on one, and `--fix` rewrites it.

## Build

`make build` compiles `./cmd` into `./tmp/telee` and stamps `cli.version` from the [`VERSION`](VERSION) file through `-ldflags`. That variable is initialised to `dev`, so a plain `go build ./cmd` produces a working binary whose `--version` reads `telee version dev`.

There is no make target for the container image, because GoReleaser builds it during a release. `dockers_v2` in [`.goreleaser.yml`](.goreleaser.yml) lays the context out as `<os>/<arch>/telee` beside `LICENSE` and `NOTICE`, and [`Dockerfile`](Dockerfile) copies `$TARGETPLATFORM/telee` out of it onto `scratch`. Images are pushed to `ghcr.io/umatare5/telee`, and a prerelease is excluded from the `latest`, `vX` and `vX.Y` tags.

The release build targets `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64` with `CGO_ENABLED=0` and `-trimpath`. A post-build hook greps `CGO_ENABLED=0` out of each binary, so cgo creeping back in fails the release rather than shipping.

## Release

To release a new version, follow these steps:

1. Rename `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md) to `## [vX.Y.Z]`, then add that version's release link at the foot and repoint the `[Unreleased]` compare link at the new tag.
2. Update the version in the [`VERSION`](VERSION) file.
3. Refresh the coverage badge — `make test-unit`, then `mkdir -p tmp && octocov badge coverage --config .octocov.yml > tmp/coverage.svg && mv tmp/coverage.svg docs/assets/coverage.svg`. The redirect goes through `tmp/` because the shell truncates its target before octocov runs, so aiming it at the badge empties the last good one whenever the profile is missing. Nothing automates the step: the shared coverage workflow checks the threshold and writes no badge.
4. Submit a pull request with all three files.

Merging that pull request is the whole release. A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/telee/actions/workflows/go-release.yml). Its `tagging` job reads `VERSION`, pushes the annotated tag `vX.Y.Z` and waits for origin to show it, and its `release` job then runs GoReleaser against that tag.

The workflow declares `on: push` alone, so there is no manual trigger and no step to perform by hand.

GoReleaser writes the release notes itself from the commit subjects, grouping them into Features, Bug Fixes, Documentation Updates and Others, and dropping `ci:`, `test:` and `release:` subjects. [`CHANGELOG.md`](CHANGELOG.md) is the hand-written counterpart and is not generated from those notes.

## Sample identities

Every transcript in this repository is device output with its identities replaced, and this section is the canon those replacements come from. It covers [`README.md`](README.md) and the reference pages under `docs/` equally.

Three kinds have a range reserved for exactly this, so nothing here is invented:

| Kind         | Range                                     | Reserved by           |
| :----------- | :---------------------------------------- | :-------------------- |
| IPv4 address | `192.0.2.0/24`                            | RFC 5737              |
| MAC address  | `00:00:5e:00:53:00` – `00:00:5e:00:53:ff` | RFC 7042 §2.1.2       |
| Domain name  | `example.internal`                        | ICANN private-use TLD |

RFC 1918 space is excluded outright, because every address in these documents sits on a `-H` a reader can paste and `192.168.0.0/16` may be live on that reader's own LAN. The MAC block sits inside `00-00-5E`, the OUI the IEEE assigned to IANA, and its first octet leaves the I/G bit clear, so no sample address can be multicast or collide with another vendor's assignment. A transcript keeps the separator and case its device prints, which is why the same address reads `00-00-5E-00-53-01` under AlliedWare.

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
2. Create a feature branch — the `no-commit-to-main` hook rejects a commit made on `main`
3. Commit your changes, using [Conventional Commits](https://www.conventionalcommits.org/) and signing off with `Signed-off-by:`
4. Add tests beside the code under test as `*_test.go`, and update the documentation alongside the code
5. Take every sample identity from [Sample identities](#sample-identities), never a value read off a device
6. Run `make lint` and `make test-unit`
7. Rebase your local changes against the `main` branch
8. Create a new Pull Request

Four workflows run on every pull request against `main`: Format and Lint pinning golangci-lint v2.13.2, Test and Build and Coverage each pinning Go 1.27.1, and CodeQL. Four more are filtered on `paths` and run only when what they watch changes. Those are actionlint on `.github/workflows/**`, govulncheck on the Go sources and the module files, and markdownlint and Link Check on `**/*.md`.

> [!NOTE]
> Coverage is a report rather than a gate here. [`.github/workflows/go-test-coverage.yml`](.github/workflows/go-test-coverage.yml) passes `coverage_threshold: 0` and [`.octocov.yml`](.octocov.yml) accepts `0%`, so no run fails on it until both are raised.
