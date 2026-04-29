# AI Assistants – Repository Instructions (umatare5/telee)

**Last Updated**: 2026-04-30

> [!IMPORTANT]
> **Read [README.md](README.md) first** to understand the project's purpose, architecture, and constraints before making any changes.

## Overview

`telee` is a **production-grade network automation CLI** (5+ years, v1.9.0) executing single commands on network devices via Telnet/SSH.

**Primary Goal**: Maintain backward compatibility while improving code quality.

**Key Constraints**:

- CLI interface must remain stable (flags, output format, exit codes)
- Shell-friendly design (stdout: device output, stderr: diagnostics)
- Pipeline-oriented (`|`, `>`, `grep` compatibility)
- Minimal dependencies (prefer standard library)

---

## Core Principles

> [!NOTE]
> These principles guide all development decisions:

**Stability Over Features**  
This is a mature tool. Prefer bug fixes and code quality improvements over new features. Any change that could break existing scripts requires explicit approval.

**Production Quality**  
All code must be lint-clean, well-tested, and follow Go idioms. Use named constants, predicate functions, and early returns. Avoid clever code.

**Security First**  
Never log credentials. Mask all secrets in output. Prefer SSH in documentation. No session transcripts unless explicitly requested.

**Minimal Diffs**  
Change only what's necessary. Preserve existing patterns and style unless there's a compelling reason to refactor.

---

## Technical Stack

**Language & Architecture**:

- Go 1.25+ with Clean Architecture pattern
- Entry: `cmd/` → Interface: `cli/` → Business: `internal/` → Libraries: `pkg/`

**Supported Platforms**:

- Cisco: IOS/IOS-XE, NX-OS, ASA, AireOS
- Juniper: SRX, SSG
- Others: YAMAHA RT, AlliedWare, IronWare

**Network Protocols**: SSH (preferred), Telnet (legacy)

**Tooling**: Follow `.golangci.yml`, `.editorconfig`, `Makefile` — no ad-hoc overrides.

---

## Development Guidelines

### Adding Code

When adding new code:

- Check existing patterns first (consistency matters)
- Use named constants for timeouts, regex patterns, retry counts
- Factor predicate helpers (`isPlatformX()`, `hasPrivilege()`)
- Comment platform-specific assumptions
- Update README platform matrix if adding support

### Modifying CLI

> [!WARNING]
> CLI changes can break user scripts. Exercise extreme caution.

When modifying the CLI interface:

- Preserve all existing flags and environment variables
- Maintain deterministic exit codes (0=success, nonzero=failure)
- Keep `--help` concise and accurate
- Explicitly document breaking changes

### Working with Platforms

Each platform has unique prompt patterns and privilege-mode sequences:

- **Cisco IOS/IOS-XE**: Enable mode requires password, prompt changes from `>` to `#`
- **Cisco NX-OS**: Similar to IOS but with different paging commands
- **Cisco ASA**: Paging not supported in user-level mode
- **Cisco AireOS**: No enable mode, WLC-specific prompts
- **Juniper SRX/SSG**: CLI mode vs configuration mode, prompts: `>` / `#`
- **YAMAHA RT**: Japanese vendor, unique command syntax

Reference existing code in `internal/application/usecases/` and `internal/infrastructure/repositories/` for implementation patterns.

### Quality Checks

Run these commands before committing:

```bash
make build      # Must succeed without errors
make lint       # Must be clean (0 issues)
make test-unit  # Must pass all tests
```

---

## Security & Privacy

> [!CAUTION]
> **Never log credentials.** This is a production tool handling sensitive data.

**Critical Security Rules**:

- Never log `TELEE_PASSWORD`, `TELEE_PRIVPASSWORD`, or any credentials
- Mask secrets when displaying (show first 4 chars only: `${TOKEN:0:4}…`)
- Use SSH in documentation examples (Telnet only with explicit warnings)
- No automatic session logging (user must explicitly enable)

**Safe Example**:

```bash
telee -H device.example.net -C "show version" --secure
```

---

## Release Process

**Version Management**: Tracked in `VERSION` file. Follow [SemVer](https://semver.org/) strictly.

- **MAJOR**: Breaking CLI changes (flag removal, output format changes)
- **MINOR**: New features, platform support
- **PATCH**: Bug fixes, documentation

**Automation**: Release via `.goreleaser.yml`. Docker images published to `ghcr.io/umatare5/telee`.
