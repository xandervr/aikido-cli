# aikido-cli API Compliance Update

**Date:** 2026-05-06
**Status:** Implemented

## Goal

Keep `aikido-cli` aligned with the current Aikido public REST API docs and make
every documented operation reachable from the CLI.

## Current API Baseline

The checked-in endpoint catalog in `internal/cli/commands/api.go` reflects the
public ReadMe OpenAPI docs updated on 2026-05-26. It contains 145 HTTP
operations across 129 documented paths.

### 2026-05-26 sync

Two endpoints added since the prior baseline:

- `GET /task_tracking/integrations` — exposed as `aikido tasks integrations`.
- `GET /users/ide/adoption` — exposed as `aikido users ide-adoption`.

The authenticated live spec remains available through:

```bash
aikido workspace introspect
```

## Compatibility Rules

- First-class commands should not send undocumented query parameters.
- If a documented endpoint has a stable, commonly used shape, expose a named
  command in its resource group.
- If a documented endpoint has a variable request body, use `endpointCommand`
  so the command supports `--query key=value`, `--body JSON`, `--body-file`,
  and `--out`.
- `aikido api get|post|put|delete <path>` is the full-coverage escape hatch for
  every documented REST operation.
- Destructive deletes must require `--confirm` and return exit code 3 when the
  flag is missing.

## Drift Fixed

- `workspace config-errors` now calls `GET /workspace/configurationErrors`.
- `issues list --severity` filters client-side instead of sending undocumented
  `filter_severity`.
- `repos list --team` filters client-side instead of sending undocumented
  `team_id`.
- `activity --from` and `--to` convert dates/timestamps to the integer
  `start`/`end` parameters documented by Aikido.

## Coverage Added

New or expanded groups include `api`, `domains`, `local-scan`,
`endpoint-protection`, `code-quality`, `access-tokens`, `bug-bounty`, and
additional documented operations under repos, issues, teams, users, containers,
clouds, apps, licenses, webhooks, compliance, pentest, and tasks.
