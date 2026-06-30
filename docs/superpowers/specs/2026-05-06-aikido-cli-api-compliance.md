# aikido-cli API Compliance Update

**Date:** 2026-05-06
**Status:** Implemented

## Goal

Keep `aikido-cli` aligned with the current Aikido public REST API docs and make
every documented operation reachable from the CLI.

## Current API Baseline

The checked-in endpoint catalog in `internal/cli/commands/api.go` reflects the
public ReadMe OpenAPI docs fetched on 2026-06-30. It contains 161 HTTP
operations across the documented public REST and OAuth OpenAPI pages.

### 2026-06-30 sync

The docs snapshot adds 16 operations since the prior baseline:

- `POST /token` — covered by `aikido auth login` / `aikido auth refresh`.
- `GET /workspace/slaSettings` — exposed as `aikido workspace sla-settings`.
- `GET /containers/{container_repo_id}/runners` — exposed as `aikido containers runners`.
- `GET /firewall/apps/{app_id}/users` — exposed as `aikido apps users`.
- `GET /report/ciScans/issueActions` — exposed as `aikido pr-checks issue-actions`.
- `GET /issues/groups/{issue_group_id}/notes` — exposed as `aikido issues notes`.
- `PUT /issues/{issue_id}/solve` — exposed as `aikido issues solve`.
- `POST /repositories/code/{code_repo_id}/labels` — exposed as `aikido repos add-label`.
- `POST /repositories/code/{code_repo_id}/labels/{label_id}` — exposed as `aikido repos update-label`.
- `DELETE /repositories/code/{code_repo_id}/labels/{label_id}` — exposed as `aikido repos remove-label --confirm`.
- `GET /endpoint-protection/devices` — exposed as `aikido endpoint-protection devices`.
- `GET /endpoint-protection/installed-packages` — exposed as `aikido endpoint-protection installed-packages`.
- `GET /endpoint-protection/permission-groups` — exposed as `aikido endpoint-protection permission-groups`.
- `GET /endpoint-protection/{ecosystem}/exceptions` — exposed as `aikido endpoint-protection exceptions`.
- `POST /endpoint-protection/{ecosystem}/exceptions` — exposed as `aikido endpoint-protection add-exception`.
- `DELETE /endpoint-protection/exceptions/{package_exception_id}` — exposed as `aikido endpoint-protection remove-exception --confirm`.

Two summaries were also realigned with the docs:

- `GET /endpoint-protection/activityLogs` — `List endpoint activity logs`.
- `GET /licenses` — `List & Search SBOM`.

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
