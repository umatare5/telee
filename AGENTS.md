# AI Assistants – Repository Instructions (umatare5/telee)

**Last Updated**: 2026-04-30

## Overview

`telee` is a production network automation CLI (5+ years, v1.9.0) that executes single commands on network devices via Telnet/SSH. This tool is actively used in production environments and stability is paramount.

**Primary Goal**: Maintain backward compatibility while improving code quality and reliability.

**Key Constraints**:

- CLI interface must remain stable (flags, output format, exit codes)
- Shell-friendly design (stdout for device output, stderr for diagnostics)
- Pipeline-oriented (`|`, `>`, `grep` must work reliably)
- Minimal dependencies (prefer standard library)

## Core Principles

**Stability Over Features**: This is a mature tool. Prefer bug fixes and code quality improvements over new features. Any change that could break existing scripts requires explicit approval.

**Production Quality**: All code must be lint-clean, well-tested, and follow Go idioms. Use named constants, predicate functions, and early returns. Avoid clever code.

**Security First**: Never log credentials. Mask secrets in all output. Prefer SSH in documentation. No session transcripts unless explicitly requested.

**Minimal Diffs**: Change only what's necessary. Preserve existing patterns and style unless there's a compelling reason to refactor.

## Technical Context

**Stack**:

- Go 1.25+ with Clean Architecture pattern
- Supports: Cisco IOS/NX-OS/ASA/AireOS, Juniper SRX/SSG, YAMAHA, AlliedWare, IronWare
- Platform-specific prompt detection and privilege-mode handling
- Network protocols: SSH (preferred), Telnet (legacy support)

**Architecture**:

```text
cmd/         → Entry point
cli/         → CLI interface (flags, version)
internal/    → Business logic (Clean Architecture)
  domain/    → Core types, constants
  application/ → Use cases per platform
  infrastructure/ → Repositories (SSH/Telnet clients)
  framework/ → Platform executor
pkg/         → Reusable libraries
```

**Tooling**: Follow `.golangci.yml`, `.editorconfig`, `Makefile`. No ad-hoc linter disables.

## Development Guidelines

**When Adding Code**:

- Check existing patterns first (consistency matters)
- Use named constants for timeouts, regex patterns, retry counts
- Factor predicate helpers (`isPlatformX()`, `hasPrivilege()`)
- Comment platform-specific assumptions
- Update README's platform matrix if adding support

**When Modifying CLI**:

- Preserve all existing flags and environment variables
- Maintain deterministic exit codes (0=success, nonzero=failure)
- Keep `--help` concise and accurate
- Document breaking changes explicitly

**When Working with Platforms**:

- Each platform has unique prompt patterns and enable-mode sequences
- Test thoroughly (or note in PR if manual testing is needed)
- Reference existing platform implementations for patterns
- Update the "Verified On" table in README

**Code Quality**:

- Ensure `make build`, `make lint`, `make test-unit` pass
- Use `go mod tidy` after dependency changes
- Place temp artifacts under `./tmp`
- No secrets in code, tests, or examples

## Security & Privacy

**Critical Rules**:

- Never log `TELEE_PASSWORD`, `TELEE_PRIVPASSWORD`, or any credentials
- Mask secrets when displaying (show first 4 chars only: `${TOKEN:0:4}…`)
- Use SSH in documentation examples (Telnet only with warnings)
- No automatic session logging (user must explicitly enable)

**Safe Example**:

```bash
telee -H device.example.net -C "show version" --secure
```

## Platform-Specific Notes

**Cisco IOS/IOS-XE**: Enable mode requires password, prompt changes from `>` to `#`

**Cisco NX-OS**: Similar to IOS but different paging commands

**Cisco ASA**: Paging not supported in user-level mode

**Cisco AireOS**: No enable mode, WLC-specific prompts

**Juniper SRX/SSG**: CLI mode vs configuration mode, `>`/`#` prompts

**YAMAHA RT**: Japanese vendor, unique command syntax

Consult existing code in `internal/application/usecases/` and `internal/infrastructure/repositories/` for platform-specific details.

## Release Process

Version is tracked in `VERSION` file. Follow SemVer strictly:

- **MAJOR**: Breaking CLI changes (flag removal, output format changes)
- **MINOR**: New features, new platform support
- **PATCH**: Bug fixes, documentation updates

Release via `.goreleaser.yml` automation. Docker images published to ghcr.io.
