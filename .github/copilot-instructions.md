# GitHub Copilot Agent Mode – Repo Instructions (umatare5/telee)

### Scope & Metadata

- **Last Updated**: 2025-08-11

- **Precedence**: Highest in this repository (see §2)

- **Goals**:

  - **Primary Goal**: Maintain and improve the **`telee`** CLI without breaking users’ workflows.
  - Keep CLI flags and outputs **clear, small, and stable**; maximize **readability/maintainability**.
  - Prefer **minimal diffs** (avoid unnecessary churn).
  - Defaults are secure where feasible (no secret logging; prefer SSH in examples).

- **Non‑Goals**:

  - Creating new binaries or sidecar tools unrelated to `telee`.
  - Introducing editor‑external, ad‑hoc lint rules or style changes.
  - Emitting or persisting secrets/test credentials.

---

## 0. Normative Keywords

- NK-001 (**MUST**) Interpret **MUST / MUST NOT / SHOULD / SHOULD NOT / MAY** per RFC 2119/8174.

## 1. Repository Purpose & Scope

- GP-001 (**MUST**) Treat this repository as a **Go-based network automation CLI** that executes single commands on remote devices over **Telnet/SSH**, with platform-specific prompt/privilege handling via the `--exec-platform` (`-x`) option.
- GP-002 (**MUST**) Preserve `telee` as a **shell-friendly, pipeline-oriented** tool (stdout for device output; errors and diagnostics to stderr). Avoid behavior that breaks common shell usage (`|`, `>`, `grep`, etc.).
- GP-003 (**SHOULD**) Keep vendor/platform support consistent with the existing matrix (e.g., `ios`, `nxos`, `asa`, `aireos`, `srx`, `ssg`, `yamaha`, etc.). If adding a platform, extend tests/docs and the verified table accordingly.

## 2. Precedence & Applicability

- GA-001 (**MUST**) When editing/generating code in this repository, Copilot **must follow** this document.
- GA-002 (**MUST**) In this repository, **this file (`copilot-instructions.md`) has the highest precedence** over any other instruction set. **On conflict, prioritize this file**.
- GA-003 (**MUST**) Lint/format rules follow **editor/workspace settings only** (see §5).
- GA-004 (**MUST**) Review behavior is defined by this file as well.

## 3. Expert Personas (for AI edits/reviews)

- EP-001 (**MUST**) Act as a **Go 1.24** expert.
- EP-002 (**MUST**) Act as a **network automation CLI** expert (terminal UX, SSH/Telnet nuances, enable-mode, prompt matching, timeouts).
- EP-003 (**SHOULD**) Be familiar with **Cisco IOS/NX-OS/ASA/AireOS** and **Juniper SRX/SSG**, Yamaha, etc., plus their typical prompts/privilege models.

## 4. Security & Privacy

- SP-001 (**MUST NOT**) Log tokens or credentials. Mask secrets in any diagnostic text (e.g., `${TOKEN:0:4}…`).
- SP-002 (**MUST**) Prefer **SSH examples** and safe defaults in documentation. **MUST NOT** silently change existing runtime defaults that users rely on (e.g., Telnet vs SSH). If proposing a default change, open a minimal PR and clearly label it as **\[MAJOR]**.
- SP-003 (**MUST**) Never persist session transcripts unless explicitly requested by a flag or documented feature.

## 5. Editor‑Driven Tooling (single source of truth)

- ED-001 (**MUST**) Follow repository settings: `.golangci.yml`, `.editorconfig`, `.markdownlint.json`, `.markdownlintignore`, `.husky.yaml`, `Makefile`, `renovate.json`.
- ED-002 (**MUST NOT**) Add flags/rules or inline disables that are not configured.
- ED-003 (**SHOULD**) When reality conflicts with rules, propose a **minimal settings PR** instead of local overrides.

## 6. Coding Principles (Basics)

- GC-001 (**MUST**) Apply **KISS/DRY** and keep code quality high.
- GC-002 (**MUST**) Avoid magic numbers; **use named constants** proactively (timeouts, prompt regex windows, retry counts, etc.).
- GC-003 (**MUST**) Factor **predicate helpers** (`is*`, `has*`, `supports*`) for readability in prompt/OS detection and mode transitions.
- GC-004 (**SHOULD**) Keep public surface small: stable CLI flags, env vars, and well-documented exit codes.

## 7. Coding Principles (Conditionals)

- CF-001 (**MUST**) Prefer predicate helpers in conditions.
- CF-002 (**MUST**) Prefer **early returns** to keep logic simple and fast.

## 8. Coding Principles (Loops)

- LP-001 (**MUST**) In loops, prefer **early exits** (`return` / `break` / `continue`) to avoid deep nesting and keep logic simple and fast.

## 9. Working Directory / Temp Files

- WD-001 (**MUST**) Place temporary artifacts (work files, coverage, test bins, etc.) **under `./tmp`**.
- WD-002 (**MUST**) Before completion, delete **zero‑byte files** (exception: keep `.keep`).

## 10. Model‑Aware Execution Workflow (when shell execution is available)

- WF-001 (**MUST**) Use `bash` for any shell actions.
- WF-002 (**MUST**) After editing Go code: prefer `make build` if present; otherwise `go build ./...` should succeed without warnings.
- WF-003 (**SHOULD**) Run `golangci-lint run` (or `make lint`) and fix until clean.
- WF-004 (**SHOULD**) If Docker image is touched, run `make goreleaser-image` (or `docker build` using the repo `Dockerfile`) until it succeeds.
- WF-005 (**MUST**) On completion: summarize actions/results into `./.copilot_reports/<prompt_title>_<YYYY-MM-DD_HH-mm-ss>.md`.

> **Note**: This repository is **WIP**; tests and some resources may not exist yet. Skip safely if targets are missing, but **do not** invent new CI flows or heavy scaffolding unless explicitly requested.

## 11. CLI Surface & UX Guarantees

- UX-001 (**MUST**) Keep `--hostname/-H`, `--command/-C`, `--secure/-s`, `--enable-mode/-e`, `--exec-platform/-x`, `--default-privilege-mode/-d`, timeout/port flags, and env vars (`TELEE_USERNAME`, `TELEE_PASSWORD`, `TELEE_PRIVPASSWORD`, etc.) **stable**.
- UX-002 (**MUST**) Ensure **deterministic exit codes**: nonzero on connection/timeout/auth/prompt-detection failures; zero on successful command execution.
- UX-003 (**SHOULD**) Maintain concise `--help` output and version string consistency (pull from `VERSION` file where applicable).
- UX-004 (**SHOULD**) Keep outputs grep-/pipe-friendly; avoid unexpected banners or prefixing unless behind a verbose flag.

## 12. Platform Matrix & Prompt Handling

- PM-001 (**MUST**) Preserve current behavior for supported platforms (`ios`, `nxos`, `asa`, `aireos`, `srx`, `ssg`, `yamaha`, `allied`, `foundry`, etc.).
- PM-002 (**MUST**) If adjusting prompt regex/enable-mode sequences, add targeted unit-style checks where feasible and update the **Verified On** table in README.
- PM-003 (**SHOULD**) Add comments near platform-specific logic explaining assumptions (e.g., paging commands not supported in ASA user-level; redundant prompts on SSG).

## 13. Dependencies & Packaging

- DP-001 (**MUST**) Prefer **standard library** first; be conservative about adding dependencies.
- DP-002 (**MUST**) Keep module tidy: `go mod tidy` after changes; avoid indirect bloat.
- DP-003 (**SHOULD**) Keep image build reproducible (`Dockerfile`, `Makefile`), using existing targets and `VERSION`-driven releases.

## 14. Tests / Quality Gate (for AI reviewers)

- QG-001 (**MUST**) Keep CI green where configured. In the absence of tests, ensure **lint + build + `telee --help`** succeed.
- QG-002 (**SHOULD**) For platform tweaks, add minimal sanity checks (e.g., regex unit tests) without introducing heavy fixtures.
- QG-003 (**MUST NOT**) Commit secrets, device IPs, or live credentials. Use synthetic data in any future tests.

## 15. Change Scope & Tone (for AI reviewers)

- CS-001 (**MUST**) Focus on the **diff**; propose broad refactors only with an explicit label (e.g., `allow-wide`).
- CS-002 (**SHOULD**) Tag comments with **\[BLOCKER] / \[MAJOR] / \[MINOR (Nit)] / \[QUESTION] / \[PRAISE]**.
- CS-003 (**SHOULD**) Structure comments as **“TL;DR → Evidence (rule/proof) → Minimal‑diff proposal.”**

## 16. Release & Versioning

- RV-001 (**MUST**) Follow the repo’s release flow based on the `VERSION` file and existing automation. Bump `VERSION` in a dedicated PR when preparing a release.
- RV-002 (**SHOULD**) Keep SemVer discipline. Any flag/behavior change that could break scripts is at least **MINOR**, often **MAJOR**.

## 17. Documentation Rules

- DOC-001 (**MUST**) Keep README examples up to date (installation, `docker run ghcr.io/umatare5/telee`, supported platforms, verified table).
- DOC-002 (**MUST**) Prefer SSH in examples; when showing Telnet, add a short caution note.
- DOC-003 (**MUST NOT**) Leak secrets in examples. Use placeholders and mask patterns.
- DOC-004 (**SHOULD**) When adding a platform or changing prompts, update the **Syntax/Usage** and **Verified On** sections accordingly.

## 18. Quick Checklist (before completion)

- QC-001 (**MUST**) Build succeeds (`make build` or `go build ./...`).
- QC-002 (**MUST**) Lint/format are clean per editor settings (no ad‑hoc flags/inline disables).
- QC-003 (**MUST**) README reflects any user-visible change (flags, examples, platform matrix, release steps).
- QC-4 (**SHOULD**) Docker image builds if touched (`make goreleaser-image`).
- QC-005 (**MUST**) Temp artifacts under `./tmp`, zero‑byte files removed, and report written to `./.copilot_reports/`.

---

### Appendix A — Known CLI Conventions (Do Not Break)

1. **Stdout is sacred** → Device command output only. Diagnostics to **stderr**.
2. **Single-shot UX** → Command executes once and exits; avoid long-running background behaviors unless explicitly enabled.
3. **Non-interactive first** → Fail fast with clear messages when prompts/privileges cannot be established.
4. **Masking** → Any echo of credentials/envs must be masked by default.

### Appendix B — Safe Examples Snippets

- SSH (preferred):

```bash
telee -H host.example.net -C "show int status" --secure
```

- Telnet with caution:

```bash
TELEE_USERNAME=user TELEE_PASSWORD='******' telee -H 192.0.2.10 -C "show ver" -x ios
```

> Keep examples short and copy/paste-able. Avoid real hostnames/IPs.
