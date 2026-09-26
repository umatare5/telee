# Security Policy

This policy covers `telee` and the container image published for it.

## Supported Versions

Only the most recent tagged release carries fixes, because no older tag gets a patch branch.
Reproduce a finding against that release before reporting it.

## Reporting a Vulnerability

Report privately through [GitHub Security Advisories](https://github.com/umatare5/telee/security/advisories/new), never through an issue or a pull request.

The response is best effort, and no reply time is promised.

## What to Include

**Redact these first.** Everyone invited to an advisory thread can read it, so none of them belongs in a report.

- The login or enable password, from `--password`, `--priv-password` or either environment variable
- The device hostname and its management address, both of which the invocation carries
- The captured output, which is a running configuration whenever the command asked for one

Then include the following.

- **Affected versions** – the release or image tag reproduced against, and the device's OS version.
- **Exec platform** – the `--exec-platform` value, because each one drives a separate session script.
- **Reproduction steps** – the command and its flags, with every value above removed.
- **Transport** – whether it reproduces over telnet, under `--secure-mode`, or both.
- **Impact** – state the exploit scenario, and what it reaches.
- **Suggested fix** – propose a remediation where you have one; this one is optional.
- **Disclosure status** – say whether it is shared elsewhere, and give your plan for sharing it.

## Exposure

The login and enable passwords are the device's own, and telnet carries both in clear text.
Keep them out of the places another person can read, in order of preference.

1. **`$TELEE_PASSWORD`, `$TELEE_PRIVPASSWORD` from `read -rs`.** Visible in the process environment.
2. **The same variables exported inline.** Visible in the shell history as well.
3. **`--password` and `--priv-password`.** Visible in the process list and the shell history.

- **Posture** – each exposure is documented, not accidental, so keep the passwords on a controlled path.
- **Transport** – telnet is the default, and `--secure-mode` moves the session to SSH, whose host key check fails closed.
- **Output** – nothing is logged, but a redirect of `show run` writes a running configuration, so treat that file as a backup.

> [!IMPORTANT]
> A leaked password is a leaked device login. Rotating it requires changing the password on the device or its AAA server.

## In Scope

The following fall within this policy.

- A credential reaching stdout, stderr or a log line, none of which this CLI masks
- Host key verification weakened by anything but an operator's own `--host-key-path`
- A session reaching a host the invocation did not name
- The published container image, because a defect in the image `ghcr.io` serves may not exist in the source

## Out of Scope

The following fall outside this policy.

- Telnet carrying the credential in clear text, which is the protocol and which `--secure-mode` answers
- A credential visible in the process list after `--password` or `--priv-password`, which [Exposure](#exposure) ranks last
- A device-side or vendor-OS defect, which belongs to that vendor's PSIRT
- A dependency advisory with no code path reachable from `./cmd`, unless you show the reachable path
- An operator's own configuration, which [Customization](README.md#customization) covers
