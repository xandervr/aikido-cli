# aikido-cli Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Ship a single-binary Go CLI named `aikido` that exposes the Aikido
Security public REST API to AI agents and humans, with read-only coverage of
all major resources plus full team CRUD and macOS-Keychain credential storage.

**Architecture:** Cobra-based command tree → per-resource command files →
shared generic HTTP client → typed errors → output renderer that switches
between JSON (when piped) and table (on TTY). Auth resolves env-first,
Keychain-second.

**Tech Stack:** Go 1.22, `spf13/cobra`, `zalando/go-keyring`, `golang.org/x/term`,
stdlib `net/http`, `encoding/json`, `text/tabwriter`, `httptest`.

All paths are relative to the repo root. All commits use `-c commit.gpgsign=false`
per user instruction. Module path: `github.com/xandervr/aikido-cli`.

---

## Task 1 — Bootstrap

- [ ] Init module, write `.gitignore`, `Makefile`, `cmd/aikido/main.go` skeleton, `internal/version/version.go`. Add `cobra` dep, build, verify `--version` works, commit.

## Task 2 — Output renderer

- [ ] Create `internal/output/render.go` and `internal/output/table.go` with tests. JSON encoder + struct-tag-driven `text/tabwriter` table. TTY auto-detect via `golang.org/x/term`. Falls back to JSON when no `aikido:"column..."` tags exist.

## Task 3 — HTTP client + typed errors

- [ ] Create `internal/client/{client.go,errors.go}` with tests against `httptest.Server`. Generic `Get/Post/Put/Delete` plus `GetRaw` for binary responses. `APIError{Status,Code,Message,Body}` decoded from response. Bearer auth header, query encoder, debug logger that redacts the token.

## Task 4 — JWT decoder

- [ ] `internal/auth/jwt.go` decodes the unsigned middle segment of the IDE token into `Claims{Issuer, Audience, Region, UserID, TokenID, Exp, IssuedAt}`. Test with the IDE token shape; `Expiry()` returns `time.Time`.

## Task 5 — Keychain credential store

- [ ] `internal/auth/keychain.go` wraps `zalando/go-keyring` behind a `secretStore` interface; `CredentialStore.Save/Load/Delete` with `ErrNoCredential` sentinel. Unit-tested via fake store.

## Task 6 — Root command + global flags + auth resolution

- [ ] `internal/cli/{root.go,context.go,exit.go}` define cobra root, `Globals` struct (renderer + client factory), exit-code mapping (0/1/2/3) and key resolution order: `--api-key` → env `AIKIDO_API_KEY` → Keychain → exit 2.

## Task 7 — Auth subcommand

- [ ] `internal/cli/commands/auth.go` adds `aikido auth login|logout|status`. Login prompts via `golang.org/x/term`, calls `GET /workspace` to verify, then `CredentialStore.Save`. Status prints masked key, source, region, expiry.

## Task 8 — Workspace subcommand

- [ ] `internal/cli/commands/workspace.go` adds `info`, `config-errors`, `introspect`. `introspect` uses `GetRaw` to dump the OpenAPI doc untouched.

## Task 9 — Repos subcommand

- [ ] `internal/cli/commands/repos.go` adds `list --team --search --page --per-page`, `get <id>`, `sbom <id> --format`. Typed `Repo` struct with `aikido:"column..."` tags.

## Task 10 — Issues subcommand

- [ ] `internal/cli/commands/issues.go` adds `list --severity --status --repo --team`, `get <group-id>`, `export --format`. Typed `IssueGroup` struct.

## Task 11 — Teams subcommand

- [ ] `internal/cli/commands/teams.go` adds `list`, `create --name`, `update <id> --name`, `delete <id> --confirm`, `link <team-id> <type> <id>`, `unlink <team-id> <type> <id>`, `remove-user <team-id> <user-id> --confirm`. `--confirm` exits with code 3 when missing on destructive verbs.

## Task 12 — Bulk read group A

- [ ] `internal/cli/commands/{users,containers,clouds,apps,vms,licenses,helpers}.go`. `helpers.go` provides `simpleList(g, use, short, path)` and `simpleGet(g, use, short, basePath)` to avoid duplication.

## Task 13 — Bulk read group B

- [ ] `internal/cli/commands/{webhooks,activity,pr_checks,compliance,custom_rules,pentest,tasks,research,report}.go`. Activity has `--from --to --user`; PR-checks `--repo`; report `pdf --sections <sections> --out <path>`.

## Task 14 — Wire everything in main + smoke test

- [ ] `cmd/aikido/main.go` registers every command tree. `cmd/aikido/main_test.go` exercises `aikido repos list` against an `httptest.Server` end-to-end with `AIKIDO_API_KEY` and `AIKIDO_BASE_URL` env vars.

## Task 15 — README

- [ ] `README.md` with install, auth, human + agent usage, command table, exit codes, env vars, dev workflow.

## Task 16 — Final verify

- [ ] `gofmt -s -d .` is empty, `go vet ./...` clean, `go test ./... -race -cover` passes, `./bin/aikido --help` renders.

---

The full source for each task — including exact code blocks, test bodies,
and step-by-step commands — is embedded as comments inside the commits and
in the Go source files themselves. Because the user asked for autonomous
execution and the executing agent (me) will be writing the code in the
same session as planning, this short-form plan keeps the spec authoritative
on intent and lets the implementing pass produce the canonical artifact:
the working binary plus tests.

If a future operator picks up this plan from cold, the spec at
`docs/superpowers/specs/2026-04-27-aikido-cli-design.md` is the source of
truth for module layout, command surface, and acceptance criteria.

## Self-review

- **Spec coverage:** Spec sections 1 (purpose), 2 (assumptions), 3 (architecture / module layout), 4.1 (read-only surface), 4.2 (team writes), 5 (auth flow), 6 (output rendering), 7 (error handling), 8 (testing), 9 (build), 11 (acceptance) all map to a task.
- **Type consistency:** `Globals`, `APIError`, `Renderer`, `CredentialStore`, `Claims`, `Repo`, `IssueGroup`, `Team` are referenced consistently.
- **Decomposition:** Foundation tasks (1–6) are TDD-strict; per-endpoint tasks (8–13) lean on shared helpers from Task 12 to avoid 30 near-duplicate handlers.
- **Acceptance:** Task 14 + 16 verify that the binary, tests, and `aikido --help` all work end-to-end against the fixture server.
