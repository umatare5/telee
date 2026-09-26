# Architecture

This document preserves the foundational design and architectural principles of telee.

## Output

One invocation opens one session, runs each command in turn and prints the transcript.

- **stdout is the device's** – the transcript leaves `internal/framework/exec.go` as a single `fmt.Fprintln(os.Stdout, …)`, so a pipe or a redirect carries device text and nothing else. It runs from the prompt before the first command to the prompt after the last.
- **stderr is everything else** – a rejected argument set, the enable-mode notice, the transport's failure line with its `[Hint]` block, and the closing `ERROR` log line
- **`--help` and `--version` go to stdout** – both exit 0, so neither reads as a failure to a caller testing the status

> [!NOTE]
> A usage fault is the one case that puts non-device text on stdout. `--hostname` and `--command` are declared `Required`, and urfave answers a missing or undefined flag by writing `Incorrect Usage:` to stderr and the whole help text to stdout before exiting 1.

## Exit Codes

One invocation reaches one device and either gets an answer or does not, so there are two codes.

| Code | Meaning                                                     |
| :--- | :---------------------------------------------------------- |
| 0    | The device answered, or `--help` / `--version` was asked    |
| 1    | A usage fault, a rejected argument set, or a failed session |

- **Nothing partial reaches stdout** – the transcript is printed only once the session returns, so `telee -H 192.0.2.1 -C "show version" --timeout 3` exits 1 having written zero bytes there
- **A rejected argument set never dials** – `checkArguments` runs inside `config.New` before a socket opens, so the device saw nothing and the wording alone identifies the fault
