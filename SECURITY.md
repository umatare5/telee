# Security Policy

## Reporting Vulnerabilities

Report security vulnerabilities via [GitHub Security Advisories](https://github.com/umatare5/telee/security/advisories/new).

## Development Guidelines

> [!CAUTION]
> This tool handles production network credentials. Never log sensitive data.

**Critical Rules**:

- Never log `TELEE_PASSWORD`, `TELEE_PRIVPASSWORD`, or credentials
- Mask secrets in output (show first 4 chars: `${TOKEN:0:4}…`)
- Use SSH in documentation (Telnet with explicit warnings)
- No automatic session logging

**Safe Example**:

```bash
export TELEE_USERNAME=admin
export TELEE_PASSWORD='<password>'
telee -H device.example.internal -C "show version" --secure
```

**Unsafe Example**:

```bash
# DO NOT hardcode credentials
telee -u admin -p password123 -H device.example.internal
```
