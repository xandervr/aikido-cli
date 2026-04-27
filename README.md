# aikido-cli

Single-binary CLI for the [Aikido Security](https://aikido.dev) public REST
API. Built so AI agents (Claude Code, etc.) can list code repositories, query
open vulnerabilities, and manage teams without bespoke wrappers around `curl`.

## Install

```bash
git clone <this repo>
cd aikido-cli
make install        # puts `aikido` on PATH (via $GOPATH/bin)
# or:
make build          # produces ./bin/aikido
```

Requires Go 1.22+.

## Authenticate

```bash
aikido auth login            # prompts for the key; stores in macOS Keychain
aikido auth status           # source, masked key, region, expiry
aikido auth logout           # remove the stored credential
```

You can also set `AIKIDO_API_KEY` in your environment — env always wins over
Keychain. This is the path Claude Code agents take.

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

Output auto-switches: tables on a TTY, JSON when piped:

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

If the host environment has `AIKIDO_API_KEY` set, no `auth login` is required.

## Commands

| Group         | Subcommands                                              |
|---------------|----------------------------------------------------------|
| auth          | login, logout, status                                    |
| workspace     | info, config-errors, introspect                          |
| repos         | list, get, sbom                                          |
| issues        | list, get, export                                        |
| teams         | list, create, update, delete, link, unlink, remove-user  |
| users         | list, get                                                |
| containers    | list, get, sbom                                          |
| clouds        | list, assets                                             |
| apps          | list                                                     |
| vms           | list, sbom                                               |
| licenses      | list                                                     |
| webhooks      | list                                                     |
| activity      | (top-level — `--from --to --user`)                       |
| pr-checks     | list                                                     |
| compliance    | soc2, nis2, iso27001                                     |
| custom-rules  | list, get                                                |
| pentest       | get, attack                                              |
| tasks         | projects, list                                           |
| research      | cve, changelog, malware-packages                         |
| cve, changelog, malware-packages | top-level shortcuts                   |
| report        | pdf                                                      |

`workspace introspect` dumps the live OpenAPI spec — useful for spotting any
endpoint not yet wired into a subcommand.

## Global flags

| Flag         | Effect                                              |
|--------------|-----------------------------------------------------|
| --json       | Force JSON output                                   |
| --table      | Force table output                                  |
| --no-color   | Disable ANSI colors                                 |
| --debug      | Log HTTP requests/responses to stderr               |
| --base-url   | Override base URL                                   |
| --api-key    | Override API key                                    |

## Environment variables

| Variable         | Effect                                |
|------------------|---------------------------------------|
| AIKIDO_API_KEY   | API token (wins over keychain)        |
| AIKIDO_BASE_URL  | Override base URL                     |
| NO_COLOR         | Disable colors (standard convention)  |

## Exit codes

- `0` — success
- `1` — API or network error
- `2` — missing or invalid auth
- `3` — usage / validation error (also: `--confirm` missing on destructive ops)

## Destructive operations

`aikido teams delete` and `aikido teams remove-user` require an explicit
`--confirm` flag. Without it the command exits with code 3 and changes
nothing.

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
