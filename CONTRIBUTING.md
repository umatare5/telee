# Contributing

Thank you for considering a contribution.

## Development

Install [`gotestsum`](https://github.com/gotestyourself/gotestsum), [`golangci-lint`](https://golangci-lint.run/docs/welcome/install/local/), [`pre-commit`](https://pre-commit.com/#install) and [`gitleaks`](https://github.com/gitleaks/gitleaks#installing), then run `make pre-commit-install`.

- **Hook order** – the branch guard, `golangci-lint`, `actionlint`, `gitleaks`, then `markdownlint-cli2`.
- **The guard carries `fail_fast`** – a commit on `main` stops there, so work on a branch.
- **Only `gitleaks` comes from `PATH`** – pre-commit builds the rest at the versions it pins.
- **The markdown hook runs `--fix`** – it rewrites files, so reach it with `make pre-commit-test`.
- **One install arms every worktree** – the hook path is shared, so it passes `--allow-missing-config`.

## Commands

`make help` prints this list together with the tools each target needs.

| Command                     | Description                                              |
| :-------------------------- | :------------------------------------------------------- |
| `make help`                 | Display available targets and requirements               |
| `make build`                | Build the binary into `./tmp/telee`                      |
| `make lint`                 | Verify the lint config, run golangci-lint, tidy `go.mod` |
| `make test-unit`            | Run unit tests with coverage using gotestsum             |
| `make test-unit-coverage`   | Generate the HTML coverage report                        |
| `make snapshot`             | Build a GoReleaser snapshot                              |
| `make clean`                | Remove the build and coverage artifacts                  |
| `make pre-commit-install`   | Install the pre-commit hooks                             |
| `make pre-commit-test`      | Run every hook across the tree                           |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                              |

## Build

`make build` stamps `cli.version` from [`VERSION`](VERSION), and a plain `go build ./cmd` leaves it at `dev`.

- **No image target** – GoReleaser builds the image during a release and pushes it to `ghcr.io/umatare5/telee`.
- **Prereleases publish no image** – a prerelease skips the exact, `latest` and `vX.Y` tags alike.

## Testing

Ship a test with the change it covers, and run the suite before every commit.

1. Add the test beside the code it covers, in the `_test` package `testpackage` enforces.
2. Use identities from [Fixture Identities](#fixture-identities), never a real device.
3. Run `make test-unit` – every package under `gotestsum`, with `-race` and a coverage profile.
4. Run `make test-unit-coverage` for the HTML report under `./coverage`.

Note the following as well.

- **Coverage floor** – [CI](.github/workflows/go-test-coverage.yml) enforces `coverage_threshold`.
- **The environment is cleared** – `make test-unit` drops `TELEE_*`; add a new one to its list in the [`Makefile`](Makefile).

## Fixture Identities

A fixture copies the shape of a real device dialogue, with every identity replaced.
Every value below is synthetic, and this section is the source for the samples in the docs and the tests.

### Reserved ranges

An IPv4 address draws on `192.0.2.0/24`, which [RFC 5737][rfc5737] reserves, so none is invented.

### Defined values

The rest have no standard to draw on, so this CLI defines them:

| Kind            | Value             |
| :-------------- | :---------------- |
| Switch name     | `sw01`            |
| Serial number   | `FOC0000X0XX`     |
| Username        | `admin`           |
| Password        | `user-password`   |
| Enable password | `enable-password` |

### Exceptions

These categories are deliberately outside the scheme.

- **A test server is local** – the `pkg/ssh` tests listen on `127.0.0.1`, because the client dials a real socket
- **A captured transcript keeps its device** – the README's examples keep the `lab*` hostnames and output

> [!IMPORTANT]
> Never paste a captured hostname, address, username or password into a fixture or a sample transcript.
> Nothing in this CLI redacts one, so a value pasted by hand reaches the tree unchanged.

## Code Style

`golangci-lint` enforces what [`.golangci.yml`](.golangci.yml) configures, and `make lint` verifies that config before running it.

## Documentation

Every fact has one page that owns it, and the other pages link to it rather than restating it.

- **Headings are pinned** – [`.markdownlint-cli2.jsonc`](.markdownlint-cli2.jsonc) sets the order via `MD043`.
- **Contracts travel** – a heading change ships with its contract in the same pull request.
- **Verbatim transcript** – [`README.md`](README.md#cli-reference) carries `--help`; usage changes update it.
- **Version reads `dev`** – `make build` stamps it, so the transcript comes from `go build`.
- **[`NOTICE`](NOTICE) tracks the module set** – a change to what the binary links updates it.
- **CI checks links** – it reaches third-party hosts. Use `lychee .` to test locally.

## Release

A release is prepared in one pull request, and merging it drafts the release.

1. Move the entries under `## [Unreleased]` in [`CHANGELOG.md`](CHANGELOG.md) to a new `## [vX.Y.Z]`.
2. List the pull requests it carries, and add that version's link at the foot.
3. Update the version in the [`VERSION`](VERSION) file.

A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/telee/actions/workflows/go-release.yml), to tag, push the images and draft the release.

- **There is no manual trigger** – the push runs it, and the weekly snapshot build tags nothing.
- **A maintainer publishes the draft** – a rerun replaces the draft until then.
- **The release links 404 until the merge** – [`lychee.toml`](lychee.toml) excludes the release-tag pattern.

## Pull Requests

Open a pull request against `main`, as a draft by default.

1. Fork the repository and create a feature branch.
2. Write [Conventional Commits](https://www.conventionalcommits.org/) and add `Signed-off-by:`.
3. Add tests and update the documentation beside the code.
4. Run `make lint` and `make test-unit`, then rebase against `main`.
5. Open the pull request.

Nothing in a commit carries a credential.

[rfc5737]: https://datatracker.ietf.org/doc/html/rfc5737
