# AI Assistants – Repository Instructions (umatare5/telee)

**Last Updated**: 2026-04-30

> [!IMPORTANT]
> **Read [README.md](README.md) first** to understand the project's purpose, architecture, and constraints before making any changes.

## Overview

`telee` is a **production-grade network automation CLI** (5+ years, v1.9.0) executing single commands on network devices via Telnet/SSH.

**Primary Goal**: Maintain backward compatibility while improving code quality.

**Key Constraints**:

- ✅ CLI interface stability (flags, output format, exit codes)
- ✅ Shell-friendly design (stdout: device output, stderr: diagnostics)
- ✅ Pipeline-oriented (`|`, `>`, `grep` compatibility)
- ✅ Minimal dependencies (prefer standard library)

---

## Core Principles

> [!NOTE]
> These principles guide all development decisions:

**🛡️ Stability Over Features**  
Prefer bug fixes and quality improvements. Breaking changes require explicit approval.

**⚙️ Production Quality**  
Lint-clean, well-tested, idiomatic Go. Use named constants, predicate functions, early returns.

**🔒 Security First**  
Never log credentials. Mask all secrets. Prefer SSH in documentation.

**📝 Minimal Diffs**  
Change only what's necessary. Preserve existing patterns unless compelling reason exists.

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

- ✅ Check existing patterns (consistency matters)
- ✅ Use named constants (timeouts, regex, retry counts)
- ✅ Factor predicate helpers (`isPlatformX()`, `hasPrivilege()`)
- ✅ Comment platform-specific assumptions
- ✅ Update README platform matrix if adding support

### Modifying CLI

> [!WARNING]
> CLI changes can break user scripts. Exercise extreme caution.

- ✅ Preserve all flags and environment variables
- ✅ Maintain deterministic exit codes (0=success, nonzero=failure)
- ✅ Keep `--help` concise and accurate
- ✅ Explicitly document breaking changes

### Working with Platforms

Each platform has unique prompt patterns and privilege-mode sequences. Reference existing code in `internal/application/usecases/` and `internal/infrastructure/repositories/`.

### Quality Checks

```bash
make build  # Must succeed
make lint   # Must be clean (0 issues)
make test-unit  # Must pass
```

---

## Security & Privacy

> [!CAUTION]
> **Never log credentials.** This is a production tool handling sensitive data.

**Critical Rules**:

- 🚫 Never log `TELEE_PASSWORD`, `TELEE_PRIVPASSWORD`, or any credentials
- 🔒 Mask secrets when displaying (first 4 chars only: `${TOKEN:0:4}…`)
- 📖 Use SSH in documentation (Telnet only with explicit warnings)
- 🚫 No automatic session logging (user must explicitly enable)

**Safe Example**:

```bash
telee -H device.example.net -C "show version" --secure
```

---

## Platform-Specific Notes

| Platform | Key Characteristics |
|----------|---------------------|
| **Cisco IOS/IOS-XE** | Enable mode requires password, prompt: `>` → `#` |
| **Cisco NX-OS** | Similar to IOS, different paging commands |
| **Cisco ASA** | Paging not supported in user-level mode |
| **Cisco AireOS** | No enable mode, WLC-specific prompts |
| **Juniper SRX/SSG** | CLI vs config mode, prompts: `>` / `#` |
| **YAMAHA RT** | Japanese vendor, unique command syntax |

> [!TIP]
> Consult existing code in `internal/application/usecases/` for implementation patterns.

---

## Release Process

**Version Management**: Tracked in `VERSION` file. Follow [SemVer](https://semver.org/) strictly.

- **MAJOR**: Breaking CLI changes (flag removal, output format changes)
- **MINOR**: New features, platform support
- **PATCH**: Bug fixes, documentation

**Automation**: Release via `.goreleaser.yml`. Docker images → `ghcr.io/umatare5/telee`.
