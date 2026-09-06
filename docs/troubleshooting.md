# Troubleshooting

A failure names its stage on the first line of stderr, and five stages can return:

```text
Incorrect Usage: Required flag "hostname" not set
2026/01/01 00:00:00 ERROR failed to validate arguments error="exec-platform is not supported"
~/.ssh/known_hosts not found. Please create it by running: ssh sw01 and accepting the key
TelnetClient was failed at spawn(). You can troubleshoot using wireshark.
TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark.
```

The first three never reach the device: the parser and the argument check run before any socket opens, and the host key file is read before the dial. The fourth dialled and failed. The fifth connected and then failed to recognise the device's prompt, which is the failure most invocations actually hit.

## The exit contract

Device output goes to stdout and nothing else does — a redirect captures the session and leaves every diagnostic on the terminal:

```bash
export TELEE_PRIVPASSWORD='<enable password>'
telee -H sw01 -C "show run" -e > sw01.cfg
```

The one exception is the argument parser, which prints the full usage block on stdout when a required flag is missing. Everything else — the validation line, the enable-mode notice, the transport banner, the underlying error, the hint block and the closing log line — is written to stderr.

Exit status is `0` on success and `1` on every failure, with no third code. The closing log line names which stage returned. `failed to validate arguments` is the argument check alone, while `command execution failed` covers the parser refusal, the host key check and the session itself.

## Refused before validation

`--hostname` and `--command` are declared required, so the argument parser refuses first and the CLI's own checks never run:

```text
Incorrect Usage: Required flag "command" not set
```

The environment satisfies the requirement as well as the flag does, so this appears only when neither `-C` nor `TELEE_COMMAND` is set. Passing an explicitly empty value clears the parser and reaches the validation refusal below instead.

## Refused during validation

The checks run in a fixed order and the first failure returns, so an invocation with two faults reports only the earlier one. The list is in that order.

- **`exec-platform is not supported`** — `-x` took a value outside the nine platform names, matched exactly and case sensitively.
- **`enable-mode and default-priv-mode cannot use at once`** — `-e` and `-d` were both set; they describe two different sessions, one that escalates and one that starts privileged.
- **`redundant-mode is not supported in this platform`** — `-r` was set on anything but `asa` or `ssg`, the two platforms whose scripts carry a failover prompt suffix.
- **`secure-mode is not supported in this platform`** — `-s` was set on `allied` or `foundry`, whose scripts exist only in a telnet form.
- **`non secure-mode is not supported in this platform`** — `-s` was omitted on `srx`, whose script exists only in an SSH form.
- **`default-privilege-mode is not supported in this platform`** — `-d` was set on anything but `asa`, `ios` or `nxos`.
- **`EnableMode must be set. Terminal length expansion in user-level is not supporting.`** — an `asa` session was asked for with neither `-e` nor `-d`; `terminal pager 0` is refused at user level, so the session would stall on the device's first pager prompt.
- **`hostname must be set`** — `-H` or `TELEE_HOSTNAME` was set to an empty string, which clears the parser's required check and fails here.
- **`command must be set`** — the same, for `-C` or `TELEE_COMMAND`.
- **`TELEE_USERNAME must be set`** — `-u` or `TELEE_USERNAME` was set to an empty string; the flag's own default is `admin`, so this cannot fire unless something overwrote it.
- **`TELEE_PASSWORD must be set`** — the same, for `-p` or `TELEE_PASSWORD`, whose default is `cisco`.
- **`TELEE_PRIVPASSWORD must be set`** — `-e` was set while `--priv-password` still held its default literal `enable`. The guard compares against that literal rather than against emptiness, so an enable password that genuinely is `enable` is indistinguishable from an unset one and must be changed on the device.

## The enable-mode notice

Four platforms have no privileged mode to escalate into, so `-e` still clears the `--priv-password` check and is then ignored and reported before the session opens:

```text
[INFO] enable-mode is ignored. It's not supported.
```

`aireos`, `allied`, `srx` and `ssg` print it. The run continues at the privilege the login granted, so a command needing more will fail on the device rather than in this CLI. The notice goes to stderr, so it does not contaminate a redirect.

## The session failed at spawn()

```text
TelnetClient was failed at spawn(). You can troubleshoot using wireshark.
dial tcp 192.0.2.1:23: connect: operation timed out
```

The banner is the transport's, the line under it is the operating system's, and `SSH was failed at spawn()` is the `--secure-mode` wording of the same stage. Nothing was authenticated and no command was sent.

The second line names the cause: `no such host` is resolution, `connection refused` is a closed port, and `operation timed out` is a filtered path. For the latter two, check that the completed port is the one the device listens on — `0` completes to 23 without `--secure-mode` and to 22 with it.

Host key failure is a separate shape. A `known_hosts` mismatch reaches this banner with its guidance block above rather than below, and a `--host-key-path` mismatch reaches it with none. A missing `known_hosts` and an unreadable `--host-key-path` are refused before any dial, so neither prints a banner at all.

## The session failed at ExpectBatch()

```text
TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark.
expect: timer expired after 2 seconds
```

The transport connected and one of the expected patterns never arrived within `--timeout`, whose value the second line repeats. The hint block printed underneath names the three causes, and the second of them is the common one.

The session script builds the expected prompt out of the `--hostname` value. `ios` waits for `<hostname>>`, `foundry` for `telnet@<hostname>>`, `allied` for `Manager <hostname>>`, `srx` for `<username>@<hostname>>`, and `ssg` for `<hostname>->`.

Dialling by IP address, or by a DNS name that differs from the device's configured hostname, therefore matches none of them. `aireos` is the only platform that does not build its prompt this way, expecting the fixed string `(Cisco Controller) >`.

Two further causes produce the same failure. A wrong `--exec-platform` waits for another vendor's login prompt. An `asa` or `ssg` device in a failover pair prints a prompt suffix — `/pri/act` and `(M)` respectively — that only `--redundant-mode` accounts for.

## Host key verification failed

Five distinct messages come from the SSH host key check, and none of them sends anything:

- **`~/.ssh/known_hosts not found`** — no `--host-key-path` was given and the file does not exist. The message carries the `ssh` line that creates it, which `ssh-keyscan` cannot on a device this old.
- **`[ERROR] Host key verification failed for <host>`** — the file exists and the key does not match or is absent from it. Four remedies follow the message, including the legacy `HostKeyAlgorithms` and `KexAlgorithms` options older IOS devices need, and the `spawn()` banner prints after them.
- **`failed to read host key file`** — `--host-key-path` named a path that does not exist or cannot be opened.
- **`failed to parse host key: ssh: no key found`** — `--host-key-path` named a file holding no key line. A `.pub` line and a `known_hosts` line both parse, and `#` comments are skipped, so the file is neither.
- **`ssh: handshake failed: ssh: host key mismatch`** — `--host-key-path` parsed, and the key it pins is not the one the device presented. It follows the `spawn()` banner without a guidance block, which belongs to the `known_hosts` path alone.

There is no flag that skips verification, by design. [`configuration.md`](configuration.md) covers `--host-key-path` in full.

## The command took far longer than --timeout

`--timeout` bounds one expect step rather than the session. The ceiling for the expect phase is the value times the script's step count, which is two on `srx` and up to seven on a telnet session using `-e`. A device that answers each prompt slowly therefore takes several multiples of the flag.

The connect stage is not bounded by it at all. Neither dial passes a deadline, so the operating system's own connect timeout governs: a non-routable address with `--timeout 3` failed after 75 s on macOS 26. A run that hangs for over a minute and then reports `operation timed out` was never affected by the flag, and raising it changes nothing.
