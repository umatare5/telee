# Documentation

Reference pages for `telee`, carrying the detail behind the [README](../README.md).

| Page                                       | Focus                                                          |
| :----------------------------------------- | :------------------------------------------------------------- |
| [`configuration.md`](configuration.md)     | Every flag, its environment variable and which one wins        |
| [`troubleshooting.md`](troubleshooting.md) | A symptom, then the condition that produced it                 |
| [`measurements.md`](measurements.md)       | Every timing taken against the lab switch, and its scripts     |
| [`../CONTRIBUTING.md`](../CONTRIBUTING.md) | The make targets, the build, the release and sample identities |
| [`../CHANGELOG.md`](../CHANGELOG.md)       | What each release carries, newest first                        |
| [`../SECURITY.md`](../SECURITY.md)         | Credential handling and where to report a vulnerability        |

## Technical Information

### Output

One invocation opens one session, runs one command and prints what the device answered.

- **stdout is the device's** — the answer leaves `internal/framework/exec.go` as a single `fmt.Fprintln(os.Stdout, …)`, so a pipe or a redirect carries device text and nothing else
- **stderr is everything else** — a rejected argument set, the enable-mode notice, the transport's failure line with its `[Hint]` block, and the closing `ERROR` log line
- **`--help` and `--version` go to stdout** — both exit 0, so neither reads as a failure to a caller testing the status

> [!NOTE]
> A usage fault is the one case that puts non-device text on stdout. `--hostname` and `--command` are declared `Required`, and urfave answers a missing or undefined flag by writing `Incorrect Usage:` to stderr and the whole help text to stdout before exiting 1.

### Exit codes

One invocation reaches one device and either gets an answer or does not, so there are two codes.

| Code | Meaning                                                     |
| :--- | :---------------------------------------------------------- |
| 0    | The device answered, or `--help` / `--version` was asked    |
| 1    | A usage fault, a rejected argument set, or a failed session |

- **Nothing partial reaches stdout** — the answer is printed only once the session returns, so `telee -H 192.0.2.1 -C "show version" --timeout 3` exits 1 having written zero bytes there
- **A rejected argument set never dials** — `checkArguments` runs inside `config.New` before a socket opens, so the device saw nothing and the wording alone identifies the fault
