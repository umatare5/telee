# Troubleshooting

A failure names its stage on the first line of stderr, and five stages can return:

```text
Incorrect Usage: Required flag "hostname" not set
2026/01/01 00:00:00 ERROR failed to validate arguments error="exec-platform is not supported"
~/.ssh/known_hosts not found. Please create it by running: ssh sw01 and accepting the key
TelnetClient was failed at spawn(). You can troubleshoot using wireshark.
TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark.
```

The first three never reach the device: the parser and the argument check run before any socket opens, and the host key file is read before the dial. The fourth dialed and failed. The fifth connected and then failed to recognize the device's prompt, which is the failure most invocations actually hit.

## The exit contract

Device output goes to stdout and nothing else does – a redirect captures the session and leaves every diagnostic on the terminal:

```bash
export TELEE_PRIVPASSWORD='<enable password>'
telee -H sw01 -C "show run" -e > sw01.cfg
```

The one exception is the argument parser, which prints the full usage block on stdout when a required flag is missing. Everything else – the validation line, the enable-mode notice, the transport banner, the underlying error, the hint block and the closing log line – is written to stderr.

Exit status is `0` on success and `1` on every failure, with no third code. The closing log line names which stage returned. `failed to validate arguments` is the argument check alone, while `command execution failed` covers the parser refusal, the host key check and the session itself.

## Refused before validation

`--hostname` and `--command` are declared required, so the argument parser refuses first and the CLI's own checks never run:

```text
Incorrect Usage: Required flag "command" not set
```

The environment satisfies the requirement as well as the flag does, so this appears only when neither `-C` nor `TELEE_COMMAND` is set. Passing an explicitly empty value clears the parser and reaches the validation refusal below instead.

## Refused during validation

The checks run in a fixed order and the first failure returns, so an invocation with two faults reports only the earlier one. The list is in that order.

- **`exec-platform is not supported`** – `-x` took a value outside the nine platform names, matched exactly and case sensitively.
- **`enable-mode and default-priv-mode cannot use at once`** – `-e` and `-d` were both set; they describe two different sessions, one that escalates and one that starts privileged.
- **`redundant-mode is not supported in this platform`** – `-r` was set on anything but `asa` or `ssg`, the two platforms whose scripts carry a failover prompt suffix.
- **`secure-mode is not supported in this platform`** – `-s` was set on `allied` or `foundry`, whose scripts exist only in a telnet form.
- **`non secure-mode is not supported in this platform`** – `-s` was omitted on `srx`, whose script exists only in an SSH form.
- **`default-privilege-mode is not supported in this platform`** – `-d` was set on anything but `asa`, `ios` or `nxos`.
- **`EnableMode must be set. Terminal length expansion in user-level is not supporting.`** – an `asa` session was asked for with neither `-e` nor `-d`, and its paging command is refused at user level.
- **`hostname must be set`** – `-H` or `TELEE_HOSTNAME` was set to an empty string, which clears the parser's required check and fails here.
- **`command must be set`** – the same, for `-C` or `TELEE_COMMAND`, and for any one `-C` given empty.
- **`command must be one line`** – a `-C` value holds a line break, which the device would run as a second command that the transcript then misattributes.
- **`TELEE_USERNAME must be set`** – `-u` or `TELEE_USERNAME` was set to an empty string; the flag's own default is `admin`, so this cannot fire unless something overwrote it.
- **`TELEE_PASSWORD must be set`** – neither `-p` nor `TELEE_PASSWORD` was set, or one of them was set to an empty string. The password has no default, so this is the refusal a run with no credentials meets before it dials.
- **`TELEE_PRIVPASSWORD must be set`** – `-e` was set while neither `--pp` nor `TELEE_PRIVPASSWORD` was, or the value was empty. The privileged password has no default either, so `enable` itself is sent as given.

## The enable-mode notice

Four platforms have no privileged mode to escalate into, so `-e` still clears the `--priv-password` check and is then ignored and reported before the session opens:

```text
[INFO] enable-mode is ignored. It's not supported.
```

`aireos`, `allied`, `srx` and `ssg` print it. The run continues at the privilege the login granted, so a command needing more will fail on the device rather than in this CLI. The notice goes to stderr, so it does not contaminate a redirect.

## The session failed at spawn()

```text
TelnetClient was failed at spawn(). You can troubleshoot using wireshark.
dial tcp 192.0.2.1:23: i/o timeout
```

The banner is the transport's, the line under it is the operating system's, and `SSH was failed at spawn()` is the `--secure-mode` wording of the same stage. That wording also covers a handshake or an authentication that outlived `--timeout`, as a slow AAA server can cause, and both end in `i/o timeout` too. A channel open that did ends in `unexpected packet in response to channel open: <nil>` instead, and a pty or shell request in `EOF`.

The second line names the cause: `no such host` is resolution, `connection refused` is a closed port, and `i/o timeout` is a filtered path, reported once `--timeout` has elapsed. For the latter two, check that the completed port is the one the device listens on – `0` completes to 23 without `--secure-mode` and to 22 with it.

Host key failure is a separate shape. A `known_hosts` refusal reaches this banner with its own message above rather than below, and a `--host-key-path` mismatch reaches it with none. A missing `known_hosts` and an unreadable `--host-key-path` are refused before any dial, so neither prints a banner at all.

## The session failed at ExpectBatch()

```text
TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark.
expect: timer expired after 2 seconds
```

The transport connected and one of the expected patterns never arrived within `--timeout` seconds of the last byte, whose value the second line repeats. A write the device does not take within `--timeout` seconds fails with the same line. A device that closes the connection first ends the step at once instead, with `expect: connection closed before a match: EOF` on the second line.

Several `-C` values share the session, so a failure on any of them ends the run with nothing on stdout, whichever answered before it. A copy of the prompt inside an answer, as IOS reprints it after `?`, can end the step early, and the run then exits 0 with a transcript that ends before the last answer does.

The hint block printed underneath names the three causes, and the second of them is the common one. The session script builds the expected prompt out of the `--hostname` value. `ios` waits for `<hostname>>`, `foundry` for `telnet@<hostname>>`, `allied` for `Manager <hostname>>`, `srx` for `<username>@<hostname>>`, and `ssg` for `<hostname>->`.

Dialing by IP address, or by a DNS name that differs from the device's configured hostname, therefore matches none of them. `aireos` is the only platform that does not build its prompt this way, expecting the fixed string `(Cisco Controller) >`.

Two further causes produce the same failure. A wrong `--exec-platform` waits for another vendor's login prompt. An `asa` or `ssg` device in a failover pair prints the suffix only `--redundant-mode` accounts for.

## Host key verification failed

Eight distinct messages come from the SSH host key check, and none of them sends anything:

- **`~/.ssh/known_hosts not found`** – no `--host-key-path` was given and the file does not exist. The message carries the `ssh` line that creates it, which `ssh-keyscan` cannot on a device this old.
- **`[ERROR] Host key verification failed for <host>`** – the file exists and holds no key for the host. Four remedies follow the message, including the legacy `HostKeyAlgorithms` and `KexAlgorithms` options older IOS devices need, and the `spawn()` banner prints after them.
- **`[ERROR] Host key for <host> does not match known_hosts: it holds TYPE at FILE:LINE …`** – another type is on record for the host, and no remedy follows because another device could present it too. The case is common, as `ssh` records the ed25519 key it prefers while this client negotiates ECDSA or RSA first. The presented type is added, or pinned with `--host-key-path`, only once `ssh-keygen -lf` gives its fingerprint for the key the device itself prints, as IOS does under `show ip ssh`.
- **`[ERROR] Host key for <host> has changed: SHA256:…`** – the file holds a key of the same type for the host and the device presented another, and the recorded line is named. The fingerprint is the device's and no remedy follows, because a changed key is what the file exists to catch. The stale line is removed only once the change is explained.
- **`[ERROR] Host key for <host> is revoked at FILE:LINE`** – an `@revoked` line in the file names the key the device presented, so no remedy follows.
- **`failed to read host key file`** – `--host-key-path` named a path that does not exist or cannot be opened.
- **`failed to parse host key: ssh: no key found`** – `--host-key-path` named a file holding no key line. A `.pub` line and a `known_hosts` line both parse, and `#` comments are skipped, so the file is neither.
- **`ssh: handshake failed: ssh: host key mismatch`** – `--host-key-path` parsed, and the key it pins is not the one the device presented. It follows the `spawn()` banner without a guidance block, which belongs to the `known_hosts` path alone.

There is no flag that skips verification, by design. [`configuration.md`](configuration.md) covers `--host-key-path` in full.
