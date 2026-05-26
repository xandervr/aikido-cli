# aikido-cli

Single-binary CLI for the [Aikido Security](https://aikido.dev) public REST
API. Built so AI agents (Claude Code, etc.) can list code repositories, query
open vulnerabilities, and manage teams without bespoke wrappers around `curl`.

## Install

```bash
git clone <this repo>
cd aikido-cli
make install                # → ~/.local/bin/aikido (no sudo)
# or:
make install-system         # → /usr/local/bin/aikido (sudo)
# or:
make build                  # only produce ./bin/aikido; install yourself
```

Requires Go 1.22+. `make install` warns if `~/.local/bin` is not on your `$PATH`.

### Shell completions

```bash
make install-completions    # writes zsh, bash, and fish completion files
```

For **zsh** (most macOS users), make sure your `~/.zshrc` contains, **before**
any `compinit` call:

```zsh
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

If completions still don't show up after running `make install-completions`,
zsh's compdump cache is stale — bust it once:

```bash
rm -f ~/.zcompdump*
exec zsh -l
```

After that, `aikido <TAB>`, `aikido issues list --<TAB>` etc. complete
subcommands and flags.

## Version

```bash
aikido version          # aikido version <version>
aikido version --json   # includes version, commit, and build date
aikido --version        # Cobra's built-in version flag
```

`make build` injects version metadata from git: `VERSION` defaults to
`git describe --tags --always --dirty`, `COMMIT` to the short commit SHA, and
`DATE` to the UTC build time. Override any of them when packaging:

```bash
make build VERSION=v1.2.3 COMMIT=abc123 DATE=2026-05-06T08:00:00Z
```

## Authenticate

The Aikido public REST API uses OAuth2 with the **client_credentials** grant.
You need a `client_id` and `client_secret` pair, generated from the Aikido
web UI:

1. Open https://app.aikido.dev → **Settings**.
2. **Integrations** → **REST API**.
3. Click **Generate token** — you'll get a **client ID** and a **client secret**.

Then run:

```bash
aikido auth login
# Aikido client ID:    <paste>
# Aikido client secret: <paste, hidden>
# ✓ Verified workspace "Focus" (id 12345)
# ✓ Stored client credentials in OS keychain
# ✓ Cached access token (valid for 59m30s)
```

The CLI does the OAuth exchange (`POST /api/oauth/token`) on your behalf,
then caches the access token under `$XDG_CACHE_HOME/aikido-cli/token.json`
(`~/Library/Caches/aikido-cli/token.json` on macOS) until it expires.

Other auth commands:

```bash
aikido auth status     # source, masked client_id, cached-token status
aikido auth refresh    # force a fresh token exchange
aikido auth logout     # delete creds + cached token
```

### Non-interactive auth (CI / agents)

Either provide the credentials in env vars and let the CLI exchange them:

```bash
export AIKIDO_CLIENT_ID=...
export AIKIDO_CLIENT_SECRET=...
aikido repos list
```

…or, if you've already exchanged a token externally, hand it directly:

```bash
export AIKIDO_ACCESS_TOKEN=...
aikido repos list
```

`AIKIDO_ACCESS_TOKEN` skips the OAuth exchange entirely.

## Examples (humans)

```bash
aikido workspace info
aikido repos list
aikido issues list --severity high
aikido issues get 12345
aikido teams list
aikido teams create --name "Platform"
aikido teams update 42 --name "Platform Eng"
aikido teams delete 42 --confirm
```

Output auto-switches: tables on a TTY, JSON when piped.

```bash
aikido issues list                # table
aikido issues list | jq '.[].id'  # JSON
aikido issues list --json         # force JSON
aikido issues list --table        # force table
```

## Examples (AI agents)

Agents pipe by default and therefore always get JSON:

```bash
aikido repos list | jq '.[] | {id, name}'
aikido issues list --severity critical | jq 'length'
aikido teams list | jq '.[] | select(.name=="Platform") | .id'
```

If `AIKIDO_CLIENT_ID` and `AIKIDO_CLIENT_SECRET` are present in the agent's
environment, no `auth login` is needed — the CLI exchanges them on demand and
caches the token.

## Commands

| Group         | Subcommands                                              |
|---------------|----------------------------------------------------------|
| auth          | login, logout, status, refresh                           |
| api           | endpoints, get, post, put, delete                        |
| workspace     | info, config-errors, introspect                          |
| repos         | list, get, sbom, scan, activate/deactivate, rules, more  |
| issues        | list, get, export, counts, issue, reachability, actions  |
| teams         | list, create, update, delete, link, unlink, add/remove-user |
| users         | list, get, ide-adoption, rights                          |
| containers    | list, get, sbom, raw-sbom, scan, registries, more        |
| clouds        | list, assets, rules, aws/azure/gcp/kubernetes, delete    |
| domains       | list, create, delete, scan, headers, subdomains          |
| apps          | list, create, get, update, delete, blocklists, events    |
| vms           | list, sbom                                               |
| licenses      | list, overwrite                                          |
| webhooks      | list, add, delete                                        |
| activity      | (top-level — `--from --to --user`)                       |
| pr-checks     | list                                                     |
| compliance    | soc2, nis2, iso27001, cis, cis-aws                       |
| custom-rules  | list, create, get, update, delete                        |
| pentest       | get, create-draft, attack                                |
| tasks         | projects, integrations, list, project-mapping, map-repos, link-task |
| local-scan    | latest                                                   |
| endpoint-protection | activity-logs                                      |
| code-quality  | findings                                                 |
| access-tokens | code-scanning                                            |
| bug-bounty    | validate-report                                          |
| research      | cve, changelog, malware-packages                         |
| cve, changelog, malware-packages | top-level shortcuts                   |
| report        | pdf                                                      |
| version       | (top-level)                                              |

`api endpoints` lists the current checked-in Aikido OpenAPI operation catalog
(145 operations from the docs snapshot used for this release). `api get|post|put|delete
<path>` is the full-coverage escape hatch for every public REST endpoint:

```bash
aikido api endpoints --search domains
aikido api get /domains --query page=0 --query per_page=20
aikido api post /domains --body '{"url":"https://example.com"}'
aikido api put /domains/42/headers --body-file headers.json
```

Named mutating commands that accept variable request bodies use the same
`--query key=value`, `--body JSON`, `--body-file path`, and `--out path` flags.

`changelog <package>` requires `--from`, `--to`, and `--language`. `report pdf`
requires `--sections` (comma-separated Aikido report sections) and accepts
`--team`. `vms sbom` defaults to `--format sbom`; use `sbom_spdx` or `csv` when
needed.

`teams link` / `teams unlink` support `repo`, `container`, `cloud`, `app`, and
`domain` resource types.

`activity --from` and `--to` accept Unix timestamps, RFC3339 timestamps, or
`YYYY-MM-DD` dates. Date-only values are converted to the integer timestamps
documented by Aikido.

`workspace introspect` dumps the live OpenAPI spec — useful for checking
whether Aikido has changed since the checked-in endpoint catalog was updated.

## Global flags

| Flag             | Effect                                              |
|------------------|-----------------------------------------------------|
| --json           | Force JSON output                                   |
| --table          | Force table output                                  |
| --no-color       | Disable ANSI colors                                 |
| --debug          | Log HTTP requests/responses to stderr               |
| --version        | Show CLI version                                    |
| --base-url       | Override API base URL                               |
| --client-id      | OAuth client ID                                     |
| --client-secret  | OAuth client secret                                 |
| --access-token   | Pre-exchanged Bearer token (skips OAuth)            |

## Environment variables

| Variable               | Effect                                                |
|------------------------|-------------------------------------------------------|
| AIKIDO_CLIENT_ID       | OAuth client ID (used to exchange for an access token)|
| AIKIDO_CLIENT_SECRET   | OAuth client secret                                   |
| AIKIDO_ACCESS_TOKEN    | Pre-exchanged Bearer token (skips OAuth)              |
| AIKIDO_BASE_URL        | Override API base URL                                 |
| NO_COLOR               | Disable colors (standard convention)                  |

## Exit codes

- `0` — success
- `1` — API or network error
- `2` — missing or invalid auth (no creds, OAuth exchange failed, 401/403)
- `3` — usage / validation error (also: `--confirm` missing on destructive ops)

## Destructive operations

Documented delete commands and `aikido teams remove-user` require an explicit
`--confirm` flag. Without it the command exits with code 3 and changes nothing.

## Development

```bash
make test        # go test ./... with race detector + coverage
make fmt         # gofmt -s -w .
make vet         # go vet ./...
make build       # build into ./bin/aikido
```

Tests use `httptest.Server` fixtures — no live API calls in CI.

## Design

See `docs/superpowers/specs/2026-04-27-aikido-cli-design.md` for the full
design document.
