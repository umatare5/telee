# Changelog

Notable changes to this project, one section per release, newest first, grouped as [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) prescribes.

This file starts at [v1.8.0], the first release built by GoReleaser from the module renamed to `github.com/umatare5/telee`. The forty-five tags before it are listed on the [releases page](https://github.com/umatare5/telee/releases) alone.

Each release from [v1.8.0] on also carries notes GoReleaser generates from the commit subjects. This file is the hand-written counterpart and is not derived from them.

## [Unreleased]

### Added

- `--command` repeats, and every value runs in the one session in the order given, so a loop over devices logs in once per device, and a value holding a line break is refused
- A `known_hosts` key of the same type that differs from the one presented is refused with the fingerprint and the recorded line, a revoked key likewise, and neither gets the onboarding steps
- A `known_hosts` record of another key type is refused the same way, naming both types, the recorded line and the presented fingerprint
- Unit tests for `pkg/expect`, `pkg/telnet` and `pkg/ssh`

### Fixed

- An IPv6 literal in `--hostname` is dialed, the address now being composed with `net.JoinHostPort`
- A device that closes the connection ends the session at once, where it previously ran out `--timeout`
- A failed session writes to stderr and exits non-zero, where it previously printed the error and the hint to stdout and returned 0
- `configor` no longer runs after the flags are assembled, so an exported `CONFIGOR_PASSWORD` can no longer replace a password given on the command line
- `--default-privilege-mode` is accepted on `asa`, whose term-length guard had rejected every session without `--enable-mode`
- `--enable-mode` on `srx` is now taken, noticed on stderr and ignored, matching `aireos`, `allied` and `ssg`, none of which has a privileged path

### Changed

- stdout carries the transcript: the prompt precedes each echoed command, where the output previously began at the echo
- The expect loop and the Telnet client are in-house on the standard library, in `pkg/expect` and `pkg/telnet`, and the SSH shell is opened directly on `x/crypto/ssh`
- `--timeout` also bounds the dial and, under `--secure-mode`, the handshake, the authentication and the shell request, where a non-routable address previously hung for the operating system's 75 s
- `--password` and `--priv-password` have no default, so a run without `TELEE_PASSWORD` stops before dialing instead of sending `cisco`, and an enable password of `enable` is accepted
- The container build moved to GoReleaser's `dockers_v2`
- The pre-commit hooks are wired through the Makefile, and `make clean` no longer reaches the worktrees under `./tmp`
- The build and the tests also run weekly ([#107])

### Removed

- `google/goexpect`, archived upstream, and `ziutek/telnet`, together with the gRPC, protobuf, genproto and goterm modules the former linked into the binary
- The vendored instruction corpus and the inert `.gemini` symlink
- The tracked coverage output

## [v1.10.2]

### Changed

- The Go toolchain moved to 1.27.1 and the lint stack was refreshed ([#105])
- The tagging and release jobs were serialized ([#106])
- The shared Renovate profile was extended and the Alpine tag pinned ([#101])
- The logo image width was set to 115px ([#104])
- Sixteen of the twenty-one commits in this release are dependency updates

### Security

- `golang.org/x/crypto` moved to v0.56.0, which Renovate flagged as a security update ([#102])

## [v1.10.1]

### Added

- `--host-key-path` / `--hkp`, sourced from `TELEE_HOSTKEYPATH`, naming a public key file for SSH host key verification
- SSH host key verification through `golang.org/x/crypto/ssh/knownhosts`, replacing the callback a code-scanning alert had flagged ([#40])
- The project logo and the centered README layout

### Changed

- The workflow set was rebuilt on the shared reusable workflows ([#84])
- The Go directive moved from 1.24 to 1.25.2
- Markdown linting moved to a `markdownlint-cli2` configuration file
- Forty-two of the ninety-nine commits in this release are dependency updates

### Fixed

- The GoReleaser build flags ([#86])

## [v1.9.0]

### Changed

- The CLI moved from `urfave/cli` v2 to v3, and its syntax was rewritten against that API
- `urfave/cli/v3` moved from v3.1.1 to v3.2.0
- The version string is read from the `VERSION` file, and GoReleaser stamps `{{ .Version }}` into the binary at release time

## [v1.8.0]

### Added

- Release builds through GoReleaser, replacing the `make release` target with a GitHub Actions workflow
- A container build, completed by the `dockers:` block in `.goreleaser.yml`
- A tag cut automatically when the `VERSION` file changes
- Renovate, through `renovate.json`

### Changed

- The Go module renamed from `telee` to `github.com/umatare5/telee`
- The Go directive moved from 1.17 to 1.24, and every module was upgraded

[Unreleased]: https://github.com/umatare5/telee/compare/v1.10.2...main
[v1.10.2]: https://github.com/umatare5/telee/releases/tag/v1.10.2
[v1.10.1]: https://github.com/umatare5/telee/releases/tag/v1.10.1
[v1.9.0]: https://github.com/umatare5/telee/releases/tag/v1.9.0
[v1.8.0]: https://github.com/umatare5/telee/releases/tag/v1.8.0
[#40]: https://github.com/umatare5/telee/pull/40
[#84]: https://github.com/umatare5/telee/pull/84
[#86]: https://github.com/umatare5/telee/pull/86
[#101]: https://github.com/umatare5/telee/pull/101
[#102]: https://github.com/umatare5/telee/pull/102
[#104]: https://github.com/umatare5/telee/pull/104
[#105]: https://github.com/umatare5/telee/pull/105
[#106]: https://github.com/umatare5/telee/pull/106
[#107]: https://github.com/umatare5/telee/pull/107
