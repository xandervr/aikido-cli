# Fig / Warp autocomplete spec

This directory ships a [Fig](https://fig.io) autocomplete spec for `aikido`.
Warp consumes Fig specs verbatim, so installing this spec gives you Warp's
floating popup with descriptions, flag suggestions, and enum-value hints
(`--severity` → `critical|high|medium|low`, `--type` → all 14 issue types,
etc.).

## Install for local development

### Option 1 — `@withfig/autocomplete-tools` (recommended)

The supported flow is via Fig's CLI:

```bash
# from this fig/ directory
npm install -g @withfig/autocomplete-tools
autocomplete dev
```

`autocomplete dev` watches `aikido.ts`, transpiles it on save, and live-loads
it into Warp / Amazon Q. The first run prints the local spec path.

### Option 2 — Manual link (fastest)

If you already have Fig or Amazon Q for CLI installed, drop the spec into
their dev-spec directory:

```bash
# Fig (legacy):    ~/.fig/autocomplete/specs
# Amazon Q (new):  ~/.local/share/amazon-q/autocomplete/specs
mkdir -p ~/.local/share/amazon-q/autocomplete/specs
ln -sf "$(pwd)/aikido.ts" ~/.local/share/amazon-q/autocomplete/specs/aikido.ts
```

Restart Warp.

## Verify it works

In Warp, type `aikido ` (with a trailing space). The popup should list every
subcommand with its description. Then try:

```
aikido issues list --severity <TAB>
aikido issues list --type <TAB>
aikido teams link 42 <TAB>            # suggests resource types
```

If the popup doesn't appear, Warp may have its built-in autocomplete
disabled (Settings → Features → "Subcommand autocomplete"). Toggle it on.

## Maintenance

The spec is hand-written. Whenever a new subcommand is added in
`internal/cli/commands/`, add a matching entry here. The schema lives at
[fig.io/docs/reference/spec](https://fig.io/docs/reference/spec).

## Upstream

Once `aikido-cli` is public, the spec can be submitted to
[withfig/autocomplete](https://github.com/withfig/autocomplete) via PR.
After it lands, every Fig / Warp / Amazon Q user gets autocomplete with no
extra setup.
