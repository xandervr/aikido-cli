# AGENTS.md

Guide for AI agents working in this repo. Keep this short and current — agents
read it on every session start.

## What this is

`aikido-cli` is a Go CLI that wraps the Aikido Security public REST API
(`https://app.aikido.dev/api/public/v1/`). Primary consumer is an AI agent
piping JSON; secondary is a human running it in a terminal.

## Build and test

```bash
make build    # → ./bin/aikido (single static binary)
make test     # go test ./... -race -cover
make fmt      # gofmt -s -w .
make vet      # go vet ./...
```

`make test` and a clean `gofmt -s -d .` (no diff) must both pass before any
commit. CI does not exist yet; do this locally.

## Commit rules

- The repo author commits without GPG signing. Use `git -c commit.gpgsign=false commit`.
- Conventional-commit prefixes: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`.
- One concept per commit. Keep messages short; let the diff carry the detail.

## Auth model — important

The Aikido public API uses **OAuth2 client_credentials**. There is no static
API key path that works.

- `aikido auth login` prompts for `client_id` + `client_secret` and stores
  them as a single JSON blob in the OS keychain (service `aikido-cli`,
  account `default`).
- On every call, the CLI exchanges the credentials at
  `POST /api/oauth/token` (HTTP Basic auth) for a Bearer access token, then
  caches it at `$XDG_CACHE_HOME/aikido-cli/token.json` (mode 0600).
- Resolution order: `--access-token`/`AIKIDO_ACCESS_TOKEN` →
  `--client-id`/`--client-secret` → `AIKIDO_CLIENT_ID`+`AIKIDO_CLIENT_SECRET`
  → cached token → keychain creds (then OAuth exchange).

For tests and CI: set `AIKIDO_ACCESS_TOKEN` to bypass the OAuth dance with a
pre-supplied bearer.

## Module layout

```
cmd/aikido/                  # entry point + CLI smoke test
internal/
  auth/                      # JWT decode, keychain creds, OAuth exchange, token cache
  client/                    # generic HTTP client (Get/Post/Put/Delete/GetRaw/Raw) + APIError
  cli/                       # cobra root, Globals, exit-code mapping
  cli/commands/              # one file per resource group; helpers.go has simpleList/simpleGet
  output/                    # Renderer: JSON when piped, table on TTY (struct-tag driven)
  version/                   # ldflags-injected version string
docs/superpowers/            # spec + plan that produced this code
```

## Conventions

- **Smart output by default**: `Renderer` switches between JSON (piped) and
  table (TTY). Force with `--json` / `--table`. Add `aikido:"column,header=Foo"`
  struct tags to enable a column in the table view.
- **Schema tolerance**: when an Aikido response shape is variable (see
  `IssueGroup` in `internal/cli/commands/issues.go`), implement
  `UnmarshalJSON` that probes alias keys via `pickStr`/`pickInt`
  (`internal/cli/commands/picks.go`), and `MarshalJSON` that returns the raw
  bytes so `--json` is never lossy.
- **Bulk read commands** use `simpleList(g, use, short, path)` /
  `simpleGet(g, use, short, basePath)` from `helpers.go`. Don't duplicate that
  logic in a new command file unless you need extra flags or a non-default
  body shape.
- **Full API coverage** lives in `api.go`: `aikido api endpoints` lists the
  checked-in documented operation catalog, and `aikido api get|post|put|delete`
  can call any public REST path with `--query`, `--body`, `--body-file`, and
  `--out`. Use `endpointCommand` for named wrappers around documented
  non-trivial operations.
- **Destructive verbs** require `--confirm` and exit with code `3` if
  missing. See `aikido teams delete` for the pattern.

## Exit codes

- `0` — success
- `1` — API/network error
- `2` — missing or invalid auth
- `3` — usage / validation error

## Adding a new endpoint

1. Find the path and shape via `aikido workspace introspect | jq '.paths."/your/path"'`
   (dumps the live OpenAPI doc the workspace serves).
2. If it's a plain GET on a list/detail, register it in the appropriate
   `internal/cli/commands/<group>.go` using `simpleList` / `simpleGet`.
3. Otherwise prefer `endpointCommand` so the wrapper gets consistent
   `--query`, JSON body, raw response, and `--confirm` behavior.
4. Add a row to the table-driven test in
   `internal/cli/commands/commands_table_test.go` to lock in the URL and method.
5. Update `documentedEndpoints` in `api.go`, `README.md`, and the Commands table.

## Things to avoid

- Don't introduce static-API-key paths. Aikido doesn't issue them for the
  public REST API.
- Don't write the Bearer token to logs or stdout. The debug logger redacts it.
- Don't paginate client-side without a `--all` flag — it surprises agents.
- Don't add backwards-compat shims for env vars that never shipped
  (e.g. `AIKIDO_API_KEY` is gone; don't bring it back).

## Reference docs

- `docs/superpowers/specs/2026-05-06-aikido-cli-api-compliance.md` — current API coverage and drift rules
- `docs/superpowers/specs/2026-04-27-aikido-cli-design.md` — original v1 design
- `docs/superpowers/plans/2026-04-27-aikido-cli.md` — implementation plan
- `README.md` — user-facing docs
- Aikido API reference: https://apidocs.aikido.dev/reference/
