# aikido-cli — Design Spec

**Date:** 2026-04-27
**Status:** Approved (pending user review of this document)
**Owner:** xander

## 1. Purpose

Build a single-binary CLI named `aikido` that exposes the Aikido Security
public REST API to humans and (primarily) to AI agents. The driving use case
is letting Claude Code agents enumerate the user's workspace — list code
repositories, query open vulnerabilities, list teams, etc. — and perform
basic team management.

Non-goals for v1:

- Mutating issues (ignore/snooze) beyond what the team-management scope needs.
- Cloud connection mutations (`POST /clouds/aws`, etc.).
- App/container/webhook CRUD.
- Multi-workspace / multi-profile support.
- A scriptable shell completion install flow.

These are deliberately deferred. The architecture leaves room to add them
later by dropping in one cobra subcommand per endpoint.

## 2. Constraints & assumptions

- **Language:** Go. Single static binary, fast cold start, easy to invoke
  from agents with no runtime to install.
- **Base URL:** `https://app.aikido.dev/api/public/v1/`. Single host (not
  region-prefixed). Overridable via `AIKIDO_BASE_URL`.
- **Auth:** `Authorization: Bearer <token>`. Token is a JWT issued by Aikido
  (the IDE token type, with claims `iss=aikido.dev`, `aud=ide.aikido`,
  `region`, `user_id`, `token_id`).
- **Region:** decoded from JWT for display only; not used to pick a base URL
  in v1.
- **Output:** primary consumer is an AI agent piping stdout. Humans should
  still get readable output when running interactively.
- **Platform:** macOS first (Keychain). Cross-platform credential storage
  comes free via `zalando/go-keyring`.

## 3. Architecture

### 3.1 Module layout

```
aikido-cli/
├── cmd/aikido/main.go              # entry; wires cobra root + subcommands
├── internal/
│   ├── auth/keychain.go            # store/load/delete API key
│   ├── auth/jwt.go                 # decode JWT claims (region, user_id)
│   ├── client/client.go            # HTTP client, base URL, retries, debug
│   ├── client/errors.go            # typed APIError + exit code mapping
│   ├── client/pagination.go        # generic page iterator
│   ├── output/render.go            # Render(any) → JSON or table by TTY
│   ├── output/table.go             # struct-tag driven table renderer
│   ├── version/version.go          # build-time version string
│   └── commands/                   # one file per resource group
│       ├── auth.go                 # login / logout / status
│       ├── workspace.go            # info / config-errors / introspect
│       ├── repos.go                # list / get / sbom
│       ├── issues.go               # list / get / export
│       ├── teams.go                # list / create / update / delete /
│       │                           # link / unlink / remove-user
│       ├── users.go                # list / get
│       ├── containers.go           # list / get / sbom
│       ├── clouds.go               # list / assets
│       ├── apps.go                 # list (Zen apps)
│       ├── vms.go                  # list / sbom
│       ├── licenses.go             # list
│       ├── webhooks.go             # list
│       ├── activity.go             # log
│       ├── pr_checks.go            # list (CI scans)
│       ├── compliance.go           # soc2 / nis2 / iso27001
│       ├── custom_rules.go         # list / get
│       ├── pentest.go              # get / attack
│       ├── tasks.go                # projects / list
│       ├── research.go             # cve / changelog / malware-packages
│       └── sbom.go                 # umbrella SBOM helpers
├── go.mod
├── Makefile                        # build / test / install / release
├── README.md                       # human + agent usage
└── docs/superpowers/specs/         # this spec lives here
```

Each command file stays small (≤150 lines). The pattern per subcommand:

1. Parse cobra flags into a typed request struct.
2. Call `client.Get[T](ctx, path, params)` (generic).
3. Hand the result to `output.Render(v)`.

No business logic in command files.

### 3.2 Dependencies

- `github.com/spf13/cobra` — command parsing.
- `github.com/spf13/pflag` — comes with cobra.
- `github.com/zalando/go-keyring` — Keychain / libsecret / Credential Mgr.
- `golang.org/x/term` — TTY detection and password prompt.
- Standard library for everything else (net/http, encoding/json, etc.).

No code generation. Wrappers are hand-written; this is faster than maintaining
an OpenAPI codegen pipeline for ~35 endpoints and gives much better flag UX.

## 4. Command surface (v1)

Global flags (apply to every command):

| Flag         | Effect                                              |
|--------------|-----------------------------------------------------|
| `--json`     | Force JSON output                                   |
| `--table`    | Force table output                                  |
| `--no-color` | Disable ANSI colors in tables / errors              |
| `--debug`    | Log HTTP request/response to stderr                 |
| `--region`   | Display-only region override                        |
| `--help`     | Cobra-generated                                     |

Environment variables:

| Variable           | Effect                                       |
|--------------------|----------------------------------------------|
| `AIKIDO_API_KEY`   | Auth token; wins over Keychain               |
| `AIKIDO_REGION`    | Display-only region override                 |
| `AIKIDO_BASE_URL`  | Base URL override (dev / proxy / future regions) |
| `NO_COLOR`         | Standard convention; disables colors         |

### 4.1 Read-only (default scope)

```
aikido auth login                 # interactive prompt; verifies via /workspace; stores in Keychain
aikido auth logout                # delete stored credential
aikido auth status                # masked key, region (from JWT), workspace name

aikido workspace info             # GET /workspace
aikido workspace config-errors    # GET /workspace/configuration-errors
aikido workspace introspect       # GET /workspace/openapispec — dumps OpenAPI doc

aikido repos list                 # --team --search --page --per-page
aikido repos get <id>
aikido repos sbom <id>            # --format json|csv|spdx (passes through)

aikido issues list                # --severity --status --repo --team --page
aikido issues get <group-id>
aikido issues export              # paginates internally; --format json|csv

aikido users list
aikido users get <id>

aikido containers list
aikido containers get <id>
aikido containers sbom <id>

aikido clouds list
aikido clouds assets              # POST /clouds/assets — read-style despite verb

aikido apps list                  # Zen

aikido vms list
aikido vms sbom <id>

aikido licenses list

aikido webhooks list

aikido activity                   # --from --to --user
aikido pr-checks list             # --repo

aikido compliance soc2
aikido compliance nis2
aikido compliance iso27001

aikido custom-rules list
aikido custom-rules get <id>

aikido pentest get <id>
aikido pentest attack <id>

aikido tasks projects
aikido tasks list <project-id>

aikido cve <cve-id>
aikido changelog <package>
aikido malware-packages

aikido report pdf                 # writes binary; --out path or stdout
```

### 4.2 Team management (writes)

```
aikido teams list
aikido teams create --name <name>
aikido teams update <id> --name <name>          # v1 supports --name only;
                                                # additional fields added once
                                                # confirmed via introspect
aikido teams delete <id>                        # --confirm required
aikido teams link <team-id> <resource-type> <resource-id>
                                                # resource-type: repo|container|cloud|vm|app|domain
aikido teams unlink <team-id> <resource-type> <resource-id>
aikido teams remove-user <team-id> <user-id>    # --confirm required
```

`--confirm` is required on destructive verbs (`delete`, `remove-user`) and
exits non-zero if missing. No `delete container` / `delete app` in v1 — out
of scope.

## 5. Auth flow

### 5.1 Login

```
$ aikido auth login
Aikido API key: ********
✓ Verified workspace "Focus" (user_id 149469, region eu)
✓ Stored in macOS Keychain (service: aikido-cli, account: default)
```

Steps:

1. Read key from `--key` flag, env var, or interactive password prompt
   (no echo).
2. Decode JWT claims to extract `region`, `user_id`, `token_id`.
3. Call `GET /workspace` to verify the key works.
4. On success, store via `keyring.Set("aikido-cli", "default", key)`.
5. On failure, print the API error and exit non-zero. Nothing stored.

### 5.2 Resolution at every command run

Order:

1. `AIKIDO_API_KEY` env var.
2. `keyring.Get("aikido-cli", "default")`.
3. Exit code `2` with message: `error: not authenticated. Run 'aikido auth
   login' or set AIKIDO_API_KEY`.

JWT decoding is best-effort and never fatal — if a non-JWT token is provided,
region falls back to "unknown" and we still attempt the request.

### 5.3 Logout

`aikido auth logout` calls `keyring.Delete("aikido-cli", "default")`. Idempotent.

### 5.4 Status

`aikido auth status` prints:

- Key source: `env` | `keychain` | `none`
- Masked key: first 8 chars + `...` + last 4
- Region (from JWT)
- Workspace name (one network call to `/workspace`)
- Token expiry (from JWT `exp`)

## 6. Output rendering

A single `output.Render(v any)` function:

```go
func Render(v any) error {
    if forcedJSON || (!forcedTable && !isTerminal(os.Stdout)) {
        return renderJSON(v)
    }
    return renderTable(v)  // falls back to JSON if v has no aikido tags
}
```

- TTY detection: `term.IsTerminal(int(os.Stdout.Fd()))`.
- Tables come from struct tags. Example:

  ```go
  type Repo struct {
      ID       int    `json:"id"        aikido:"column,header=ID"`
      Name     string `json:"name"      aikido:"column,header=Repo"`
      Severity string `json:"severity"  aikido:"column,header=Severity"`
      Internal string `json:"internal_field"` // no aikido tag → omitted from table
  }
  ```
- Detail commands (`get <id>`) often return deeply nested data; the table
  renderer falls back to pretty-printed JSON for these.
- `--json` / `--table` flags override TTY detection.
- Tables are rendered without external deps using `text/tabwriter`.
- Colors are off when `NO_COLOR` is set, when `--no-color` is passed, or when
  stdout is not a TTY.

## 7. Error handling

```go
type APIError struct {
    Status  int    // HTTP status
    Code    string // Aikido error code if present
    Message string
    Body    []byte // raw response for --debug
}
```

Mapping:

| Condition                          | Exit code | Stderr format                  |
|------------------------------------|-----------|--------------------------------|
| Success                            | 0         | —                              |
| API/network error                  | 1         | JSON when piped, text on TTY   |
| Missing/invalid auth (`401`/none)  | 2         | JSON when piped, text on TTY   |
| Usage / validation error           | 3         | Cobra default + hint           |
| `--confirm` missing for destructive | 3        | Hint to add `--confirm`        |

`--debug` logs full request line, headers (Authorization redacted), and
response body to stderr.

Retries: no automatic retries in v1 — fail fast for AI agents. We can add an
exponential-backoff retry on 429/5xx later behind a `--retry` flag.

## 8. Testing strategy

Target ≥80% coverage per the project's testing rule.

- **Unit tests**
  - `auth/jwt.go` — claim decoding, malformed token handling.
  - `auth/keychain.go` — interface mocked; happy path + missing-key path.
  - `client/errors.go` — error decoding for known Aikido error shapes.
  - `client/pagination.go` — generic page iterator.
  - `output/render.go` — JSON snapshot, table snapshot for representative
    structs, TTY auto-switching, `--json`/`--table` override behavior.
- **Integration tests**
  - One per resource group, hitting an `httptest.Server` that serves
    fixtures from `internal/client/testdata/`.
  - Fixtures committed in-repo. Captured by hand from the live API once,
    not regenerated on every test run.
- **CLI smoke tests**
  - `cmd/aikido/main_test.go` runs the binary against the same fixture
    server and asserts on stdout / exit codes.
- **No live API tests in CI.** The token is workspace-specific and would
  leak. Live verification stays manual.

## 9. Build & release

- `make build` — `go build -ldflags "-X .../version.Version=..." -o bin/aikido ./cmd/aikido`
- `make install` — `go install ./cmd/aikido` (puts `aikido` on PATH).
- `make test` — `go test ./...` plus coverage gate.
- `make release` — `goreleaser` cross-compile for darwin-amd64, darwin-arm64,
  linux-amd64, linux-arm64. Out of scope for v1 — `make install` only.
- Go version pinned via `go 1.22` in `go.mod` (any 1.22+ works).

## 10. Open questions / deferred decisions

These are explicitly out of v1 scope. Listed so we don't lose them:

- **Issue mutations** (`POST snoozeissuegroup`, `POST unsnoozeissuegroup`,
  ignore/unignore) — useful for triage workflows.
- **Webhook CRUD** — `POST /webhooks`, `DELETE /webhooks/{id}`.
- **Cloud connection** — `POST /clouds/aws`, GCP, Azure variants.
- **App/container deletion** — destructive, defer until requested.
- **User-rights changes** — `PUT /users/{id}/rights`.
- **Custom rule creation** — `POST /custom-rules`.
- **Domain CRUD + scan triggers** — `POST /domains`, `POST /domains/{id}/scan`.
- **Multi-profile support** — `--profile work` switching keychain entries.
- **Shell completions** — `aikido completion bash|zsh|fish`.

Each of these is a single cobra subcommand + client wrapper to add later.

## 11. Acceptance criteria

v1 is done when:

- `aikido auth login` prompts for, verifies, and stores the API key in the
  macOS Keychain.
- `aikido repos list` returns the workspace's repos as JSON when piped, and
  as a readable table when run interactively.
- `aikido issues list` returns vulnerabilities/issue groups with severity and
  repo filters working.
- `aikido teams create --name <name>` creates a team and returns its ID.
- `aikido teams update`, `delete`, `link`, `unlink`, `remove-user` all work
  end-to-end against the live API.
- `aikido workspace introspect` dumps the OpenAPI spec.
- All unit + integration tests pass with ≥80% coverage.
- README documents both human usage and an "AI agent usage" section with
  example prompts.
