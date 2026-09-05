# Contributing

Thank you for considering a contribution.

## Commands

The following `make` commands are available for development and testing:

| Command                     | Description                                  |
| :-------------------------- | :------------------------------------------- |
| `make help`                 | Display available targets and requirements   |
| `make build`                | Build the binary to `./tmp/telee`            |
| `make lint`                 | Run golangci-lint, then tidy `go.mod`        |
| `make test-unit`            | Run unit tests with coverage using gotestsum |
| `make test-unit-coverage`   | Generate the HTML coverage report            |
| `make clean`                | Remove build artifacts and backup files      |
| `make pre-commit-install`   | Install the pre-commit hooks                 |
| `make pre-commit-test`      | Run every hook across the tree               |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                  |

`make clean` names the binary, `tmp/dist` and `coverage/` one by one rather than removing `./tmp`, because `./tmp` is also where `git wt` puts its worktrees.

`make pre-commit-install` passes `--allow-missing-config`, because the hook path is the shared git common directory. A hook installed from one worktree fires in every other one and on the `main` checkout, where `.pre-commit-config.yaml` may not exist.

Six hooks run on a commit — `no-commit-to-main`, `golangci-lint-full`, `gofmt`, `actionlint`, `gitleaks-system` and `markdownlint-cli2`. The `golangci-lint` `rev` in [`.pre-commit-config.yaml`](.pre-commit-config.yaml) has to match `golangci_lint_version` in [`.github/workflows/go-test-fmt.yml`](.github/workflows/go-test-fmt.yml). Renovate updates the two through different datasources, so nothing forces them onto the same pull request.

`.golangci.yml` declares no `formatters` block, so `golangci-lint run --fix` rewrites logic and never formatting. The separate `gofmt` hook is what keeps the tree formatted.

## Build

`make build` compiles `./cmd` into `./tmp/telee` and stamps `cli.version` from the [`VERSION`](VERSION) file through `-ldflags`. That variable is initialised to `dev`, so a plain `go build ./cmd` produces a working binary whose `--version` reads `telee version dev`.

There is no make target for the container image, because GoReleaser builds it during a release. `dockers_v2` in [`.goreleaser.yml`](.goreleaser.yml) lays the context out as `<os>/<arch>/telee` beside `LICENSE` and `NOTICE`, and [`Dockerfile`](Dockerfile) copies `$TARGETPLATFORM/telee` out of it onto `scratch`. Images are pushed to `ghcr.io/umatare5/telee`, and a prerelease is excluded from the `latest`, `vX` and `vX.Y` tags.

The release build targets `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64` with `CGO_ENABLED=0` and `-trimpath`. A post-build hook greps `CGO_ENABLED=0` out of each binary, so cgo creeping back in fails the release rather than shipping.

## Release

To release a new version, follow these steps:

1. Rename `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md) to `## [vX.Y.Z]`, then add that version's release link at the foot and repoint the `[Unreleased]` compare link at the new tag.
2. Update the version in the [`VERSION`](VERSION) file.
3. Refresh the coverage badge — `make test-unit`, then `octocov badge coverage --config .octocov.yml > docs/assets/coverage.svg`. Nothing automates it: the shared coverage workflow checks the threshold and writes no badge.
4. Submit a pull request with all three files.

Merging that pull request is the whole release. A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/telee/actions/workflows/go-release.yml). Its `tagging` job reads `VERSION`, pushes the annotated tag `vX.Y.Z` and waits for origin to show it, and its `release` job then runs GoReleaser against that tag.

The workflow declares `on: push` alone, so there is no manual trigger and no step to perform by hand.

GoReleaser writes the release notes itself from the commit subjects, grouping them into Features, Bug Fixes, Documentation Updates and Others, and dropping `ci:`, `test:` and `release:` subjects. [`CHANGELOG.md`](CHANGELOG.md) is the hand-written counterpart and is not generated from those notes.

## Pull requests

1. [Fork](https://github.com/umatare5/telee/fork) the repository
2. Create a feature branch — the `no-commit-to-main` hook rejects a commit made on `main`
3. Commit your changes, using [Conventional Commits](https://www.conventionalcommits.org/) and signing off with `Signed-off-by:`
4. Add tests beside the code under test as `*_test.go`, and update the documentation alongside the code
5. Run `make lint` and `make test-unit`
6. Rebase your local changes against the `main` branch
7. Create a new Pull Request

Four workflows run on every pull request against `main`: Format and Lint pinning golangci-lint v2.13.2, Test and Build and Coverage each pinning Go 1.27.1, and CodeQL. A fifth, actionlint, is filtered on `paths` and runs only when a file under `.github/workflows/` changed.

> [!NOTE]
> Coverage is a report rather than a gate here. [`.github/workflows/go-test-coverage.yml`](.github/workflows/go-test-coverage.yml) passes `coverage_threshold: 0` and [`.octocov.yml`](.octocov.yml) accepts `0%`, so no run fails on it until both are raised.
