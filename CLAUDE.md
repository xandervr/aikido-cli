# CLAUDE.md

Project context for Claude Code.

The canonical agent guide for this repo is **`AGENTS.md`** — read that first.
This file adds Claude-Code-specific notes that aren't relevant to other
agents.

## Quick orientation

- Single-binary Go CLI wrapping the Aikido Security public REST API.
- Primary consumer of `aikido` is an AI agent piping JSON.
- See `AGENTS.md` for build/test/commit rules and the auth model.

## Spec and plan

This codebase was generated via the Superpowers brainstorming → writing-plans
→ executing-plans workflow:

- Spec: `docs/superpowers/specs/2026-04-27-aikido-cli-design.md`
- Plan: `docs/superpowers/plans/2026-04-27-aikido-cli.md`

When making non-trivial changes, update the spec first if intent shifts, then
the code. The plan is historical.

## Hooks and gates

This repo has no project-level hooks. The user-level fact-forcing gates from
the operator's `~/.claude` config (gateguard, destructive-command facts) apply
to any work in this directory. Provide the required facts before each Write
or destructive Bash.

## Skills that apply

- `superpowers:writing-plans`, `superpowers:executing-plans`,
  `superpowers:brainstorming` — used to produce the existing spec and plan.
- Reusing them for new features is preferred to ad-hoc implementation.

## Commit attribution

The repo author asked for commits **without GPG signing or AI co-author
trailers**. Use:

```bash
git -c commit.gpgsign=false commit -m "..."
```

No `Co-Authored-By: Claude` line. (Their global `~/.claude/settings.json`
disables that attribution; respect it here too.)

## Verifying secrets locally

```bash
security find-generic-password -s aikido-cli -a default     # metadata only
ls -la ~/Library/Caches/aikido-cli/token.json                # cached access token (0600)
./bin/aikido auth status                                     # source, expiry
```

The keychain value is a JSON blob `{"client_id":"...","client_secret":"..."}`.
The cached access token is a separate, short-lived file.
