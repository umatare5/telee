# Configuration

Every setting this CLI reads, in precedence order: a flag, then the environment, then the built-in default.

One invocation runs one command on one device, so there are no subcommands, no configuration file and no per-command flag set. The thirteen flags below are the whole surface.

## Flags

`--hostname` and `--command` are required; everything else carries a default the table names.

| Flag                             | Default  | Sets                              |
| :------------------------------- | :------- | :-------------------------------- |
| `--hostname`, `-H`               | —        | Target address and prompt literal |
| `--command`, `-C`                | —        | The command the session runs      |
| `--exec-platform`, `-x`          | `ios`    | The session script to drive       |
| `--port`, `-P`                   | `0`      | TCP port, `0` meaning complete it |
| `--timeout`, `-t`                | `5`      | Seconds allowed per expect step   |
| `--secure-mode`, `-s`            | `false`  | SSH in place of telnet            |
| `--enable-mode`, `-e`            | `false`  | Escalation to privileged EXEC     |
| `--default-privilege-mode`, `-d` | `false`  | Login is already privileged       |
| `--redundant-mode`, `-r`         | `false`  | Failover suffix on the prompt     |
| `--username`, `-u`               | `admin`  | Login user                        |
| `--password`, `-p`               | `cisco`  | Login password                    |
| `--priv-password`, `--pp`        | `enable` | Password the `enable` step sends  |
| `--host-key-path`, `--hkp`       | —        | SSH host key file, in wire format |

Three boolean flags carry long aliases as well: `--ena` and `--enable` for `--enable-mode`, `--redundant` for `--redundant-mode`, and `--sec` and `--secure` for `--secure-mode`.

What the flags that hold a mechanism actually do:

- **`--port` at `0` is completed, not defaulted.** Zero becomes 22 under `--secure-mode` and 23 without it, and any non-zero value is dialled as given with no validation.
- **`--timeout` bounds one expect step, not the session.** `ExpectBatch` applies the value per step, so the expect phase's ceiling is the value times the script's step count — two on `srx`, up to seven on a telnet session using `--enable-mode`.
- **Neither dial carries a deadline, so `--timeout` does not bound the connect.** A non-routable address with `--timeout 3` failed after 75 s on macOS 26, which is the operating system's connect timeout rather than the flag.
- **`--hostname` is also the prompt pattern.** Eight of the nine scripts expect the value verbatim inside the device prompt, so anything but the device's own hostname matches nothing; `aireos` expects `(Cisco Controller) >` instead.
- **`--priv-password` left at `enable` counts as unset.** The guard compares the value against that default literal, so `--enable-mode` without an explicit password is refused rather than sent.
- **`--host-key-path` takes a wire-format key blob.** `ssh.ParsePublicKey` reads the RFC 4253 §6.6 encoding, so a `.pub` line or a `known_hosts` line is refused with `failed to parse host key: ssh: short read`.
- **`--redundant-mode` only widens the expected prompt.** It appends `/pri/act` on `asa` and `(M)` on `ssg`, and every other platform refuses the flag outright.

Which flags a platform takes is decided by the platform, not by the operator:

| `-x` value | Product                  | Escalation             | Transport   |
| :--------- | :----------------------- | :--------------------- | :---------- |
| `aireos`   | Cisco AireOS             | None                   | Telnet, SSH |
| `allied`   | AlliedTelesis AlliedWare | None                   | Telnet      |
| `asa`      | Cisco ASA Software       | `-e` or `-d`, required | Telnet, SSH |
| `foundry`  | Brocade IronWare         | `-e`                   | Telnet      |
| `ios`      | Cisco IOS/IOS-XE         | `-e` or `-d`           | Telnet, SSH |
| `nxos`     | Cisco NX-OS              | `-e` or `-d`           | Telnet, SSH |
| `srx`      | Juniper JunOS            | None                   | SSH         |
| `ssg`      | Juniper ScreenOS         | None                   | Telnet, SSH |
| `yamaha`   | YAMAHA RT                | `-e`                   | Telnet, SSH |

A platform with no escalation ignores `-e` and prints a notice on stderr, though the `--priv-password` check still runs ahead of that notice, and it refuses `-d` outright. A platform listed for one transport refuses the other: `-s` is rejected on `allied` and `foundry`, and omitting it is rejected on `srx`.

The `asa` entry is required because paging is disabled with `terminal pager 0`, which an unprivileged ASA session will not accept. [`troubleshooting.md`](troubleshooting.md) indexes each refusal by its message.

## Environment Variables

Six variables reach the flags they name, and the CLI reads no others:

| Variable             | Flag              | Default  |
| :------------------- | :---------------- | :------- |
| `TELEE_HOSTNAME`     | `--hostname`      | —        |
| `TELEE_COMMAND`      | `--command`       | —        |
| `TELEE_USERNAME`     | `--username`      | `admin`  |
| `TELEE_PASSWORD`     | `--password`      | `cisco`  |
| `TELEE_PRIVPASSWORD` | `--priv-password` | `enable` |
| `TELEE_HOSTKEYPATH`  | `--host-key-path` | —        |

The other seven flags have no environment source: `--port`, `--timeout`, `--exec-platform` and the four mode booleans are set on the command line or left at their defaults.

> [!IMPORTANT]
> `configor` was removed from this module, so a `CONFIGOR_`-prefixed variable no longer overrides an explicit flag — nor anything else. The six names above are read directly by the flag definitions and are the complete set.

## Precedence

A flag beats the environment, the environment beats the default, and nothing else participates:

```bash
export TELEE_PRIVPASSWORD='<enable password>'
telee -H sw01 -C "show run" -e            # uses the variable
telee -H sw01 -C "show run" -e --pp other # uses "other"
```

The default is applied first and each later source overwrites it. An exported empty string therefore differs from an unset variable — `TELEE_USERNAME=''` replaces the `admin` default and the run stops at `TELEE_USERNAME must be set`.

`--hostname` and `--command` are the exception to that pattern, because both are marked required. Leaving the flag and the variable both unset ends at the argument parser with a usage message, before any of the validation in this page runs.

## Hardening

Keep the credentials out of the argument list. A value passed on `--password` or `--priv-password` is visible in the process list to every other process on the host, and lands in the shell history. The environment is read only by what can already read the process environment, which is the narrower exposure:

```bash
export TELEE_USERNAME=operator
export TELEE_PASSWORD='<password>'
export TELEE_PRIVPASSWORD='<enable password>'
telee -H sw01 -C "show version" --secure --enable
```

Host key verification fails closed and has no bypass. With `--host-key-path` unset the session verifies against `~/.ssh/known_hosts`; a missing file and a key that does not match are both errors that send nothing. No flag and no environment variable disables the check.

Device output goes to stdout, joined there only by the usage block a missing required flag prints, so a redirect captures what the device printed. For `show run` that is a running configuration carrying interface descriptions, SNMP communities and encrypted secrets — treat the target file as the configuration backup it is.

Nothing else is written anywhere. The session is never logged to a file, and the only other stream is stderr, which carries diagnostics that [`troubleshooting.md`](troubleshooting.md) indexes. Reporting policy — what is in scope and how to send it — lives in [`SECURITY.md`](../SECURITY.md).

> [!CAUTION]
> Telnet is the default transport. Without `--secure-mode` the login password, the enable password and the whole session cross the network in clear text, and port `0` completes to 23 rather than 22.
