# Security Policy

## Supported Versions

Only the most recent tagged release carries fixes — reproduce a finding against it before reporting.

Releases are tags cut from `main` with no maintenance branch behind them, so a fix reaches an operator only in the next tag. The [releases page](https://github.com/umatare5/telee/releases) names the current one.

## Reporting a Vulnerability

Private vulnerability reporting is **not enabled** on this repository, so GitHub's "Report a vulnerability" form is unavailable to anyone but a maintainer and there is no published security address. Two steps reach a private channel:

1. **Open an issue** on the [issue tracker](https://github.com/umatare5/telee/issues) naming the affected version, the `--exec-platform` value and the class of the fault, and asking for a private channel.
2. **Wait for the draft advisory.** A maintainer can open one and add a reporter to it, and that is where the reproduction, the captured output and the impact belong.

The response is best effort, with no promised window.

> [!WARNING]
> An issue is world readable from the moment it is filed. Keep the reproduction, the device output and every credential out of it until the advisory exists.

## What to Include

**Redact these first.** None of them belongs in a report, public or private.

- The login or enable password, from `--password`, `--priv-password` or either environment variable
- The device hostname and its management address, both of which the invocation carries
- The captured output, which is a running configuration whenever the command asked for one

Then include the following:

- **Affected version** (required): The `telee --version` string, and the device OS version it was seen on
- **Exec platform** (required): The `--exec-platform` value, because each one drives a separate session script
- **Invocation** (required): The flags, with every value above removed
- **Transport** (required): Whether it reproduces over telnet, under `--secure-mode`, or both
- **Impact assessment** (required): What the fault reaches, and from where
- **Suggested fix** (optional): Proposed remediation, if any
- **Disclosure status** (required): Whether it is shared elsewhere, and the plan for sharing it

## Scope

In scope:

- A credential reaching stdout, stderr or a log line, none of which this CLI masks
- Host key verification weakened by anything but an operator's own `--host-key-path`
- A session reaching a host the invocation did not name
- The published container image, `ghcr.io/umatare5/telee`

Out of scope:

- Telnet carrying the credential in clear text, which is the protocol and which `--secure-mode` answers
- A credential visible in the process list after `--password` or `--priv-password`, whose cost [`docs/configuration.md`](docs/configuration.md) sets out
- A device-side or vendor-OS defect, which belongs to that vendor's PSIRT
- A dependency advisory with no path reachable from `./cmd` — show the path, or a `govulncheck` finding
- An operator's own configuration, which [`docs/configuration.md`](docs/configuration.md) covers
