// Fig / Warp autocomplete spec for `aikido`.
//
// Schema reference: https://fig.io/docs/reference/spec
// Warp consumes Fig specs verbatim. To install locally for Warp/Fig dev,
// see ./README.md in this directory.

const severities: Fig.Suggestion[] = [
  { name: "critical", icon: "🔴" },
  { name: "high", icon: "🟠" },
  { name: "medium", icon: "🟡" },
  { name: "low", icon: "🟢" },
];

// Aikido's actual status values, observed from the live API:
// task_open, task_in_progress, task_closed, task_done, todo, new, ignored, snoozed.
// The exact set may evolve; the suggestions are convenience, not constraints.
const statuses: Fig.Suggestion[] = [
  { name: "task_open" },
  { name: "task_in_progress" },
  { name: "task_closed" },
  { name: "task_done" },
  { name: "todo" },
  { name: "new" },
  { name: "ignored" },
  { name: "snoozed" },
];

const issueTypes: Fig.Suggestion[] = [
  { name: "open_source", description: "Open-source dependency vuln (SCA)" },
  { name: "leaked_secret", description: "Secret committed to source" },
  { name: "sast", description: "Static analysis finding" },
  { name: "iac", description: "Infrastructure as Code" },
  { name: "cloud", description: "Cloud posture (CSPM)" },
  { name: "docker_container", description: "Container image vuln" },
  { name: "cloud_instance" },
  { name: "surface_monitoring", description: "DAST" },
  { name: "malware", description: "Malicious package" },
  { name: "eol", description: "End of life" },
  { name: "mobile" },
  { name: "scm_security" },
  { name: "ai_pentest" },
  { name: "license" },
];

const resourceTypes: Fig.Suggestion[] = [
  { name: "repo" },
  { name: "container" },
  { name: "cloud" },
  { name: "app" },
  { name: "domain" },
];

const complianceFrameworks: Fig.Suggestion[] = [
  { name: "soc2" },
  { name: "nis2" },
  { name: "iso27001" },
];

const shells: Fig.Suggestion[] = [
  { name: "bash" },
  { name: "zsh" },
  { name: "fish" },
  { name: "powershell" },
];

const globalOptions: Fig.Option[] = [
  { name: "--json", description: "Force JSON output" },
  { name: "--table", description: "Force table output" },
  { name: "--no-color", description: "Disable ANSI colors" },
  { name: "--debug", description: "Log HTTP requests/responses to stderr" },
  {
    name: "--base-url",
    description: "Override API base URL (env: AIKIDO_BASE_URL)",
    args: { name: "url" },
  },
  {
    name: "--client-id",
    description: "OAuth client ID (env: AIKIDO_CLIENT_ID)",
    args: { name: "id" },
    isPersistent: true,
  },
  {
    name: "--client-secret",
    description: "OAuth client secret (env: AIKIDO_CLIENT_SECRET)",
    args: { name: "secret" },
    isPersistent: true,
  },
  {
    name: "--access-token",
    description: "Pre-exchanged Bearer token (env: AIKIDO_ACCESS_TOKEN)",
    args: { name: "token" },
    isPersistent: true,
  },
  { name: ["-h", "--help"], description: "Help for the current command" },
];

const pageOptions: Fig.Option[] = [
  { name: "--page", description: "Page (0-indexed)", args: { name: "n" } },
  { name: "--per-page", description: "Page size", args: { name: "n" } },
];

const idArg: Fig.Arg = { name: "id" };

const changelogOptions: Fig.Option[] = [
  {
    name: "--from",
    description: "Current package version",
    args: { name: "version" },
    isRequired: true,
  },
  {
    name: "--to",
    description: "Target package version",
    args: { name: "version" },
    isRequired: true,
  },
  {
    name: "--language",
    description: "Package language",
    args: { name: "language", suggestions: ["JS", "PY", "GO", ".NET", "Java", "Scala", "Kotlin"] },
    isRequired: true,
  },
];

const completion: Fig.Spec = {
  name: "aikido",
  description: "Aikido Security CLI — wraps the public REST API",
  subcommands: [
    // ─── auth ────────────────────────────────────────────────────────────
    {
      name: "auth",
      description: "Manage Aikido API credentials",
      subcommands: [
        {
          name: "login",
          description: "Verify OAuth credentials and store them in the OS keychain",
          options: [
            { name: "--client-id", args: { name: "id" } },
            { name: "--client-secret", args: { name: "secret" } },
          ],
        },
        { name: "logout", description: "Delete stored credentials and cached access token" },
        { name: "status", description: "Show current authentication state" },
        { name: "refresh", description: "Force a fresh OAuth token exchange" },
      ],
    },

    // ─── workspace ───────────────────────────────────────────────────────
    {
      name: "workspace",
      description: "Workspace info and configuration",
      subcommands: [
        { name: "info", description: "Show workspace summary" },
        { name: "config-errors", description: "List workspace configuration errors" },
        { name: "introspect", description: "Dump the OpenAPI spec from the workspace" },
      ],
    },

    // ─── repos ───────────────────────────────────────────────────────────
    {
      name: "repos",
      description: "Code repositories",
      subcommands: [
        {
          name: "list",
          description: "List code repositories",
          options: [
            { name: "--team", description: "Filter by team ID", args: { name: "id" } },
            { name: "--search", description: "Name search", args: { name: "query" } },
            ...pageOptions,
          ],
        },
        { name: "get", description: "Get a single code repository", args: idArg },
        {
          name: "sbom",
          description: "Export the SBOM for a code repo",
          args: idArg,
          options: [
            {
              name: "--format",
              description: "Output format passthrough",
              args: { name: "format", suggestions: ["json", "csv", "spdx"] },
            },
          ],
        },
      ],
    },

    // ─── issues ──────────────────────────────────────────────────────────
    {
      name: "issues",
      description: "Open issues / vulnerabilities",
      subcommands: [
        {
          name: "list",
          description: "List open issue groups (vulnerabilities)",
          options: [
            {
              name: "--severity",
              description: "Filter by severity",
              args: { name: "level", suggestions: severities },
            },
            {
              name: "--status",
              description: "Filter by group status",
              args: { name: "status", suggestions: statuses },
            },
            {
              name: "--type",
              description: "Filter by issue type",
              args: { name: "type", suggestions: issueTypes },
            },
            { name: "--repo", description: "Filter by code repo ID", args: { name: "id" } },
            { name: "--team", description: "Filter by team ID", args: { name: "id" } },
            ...pageOptions,
          ],
        },
        {
          name: "get",
          description: "Get details for an issue group",
          args: { name: "group-id" },
        },
        {
          name: "export",
          description: "Export all issues (paginates server-side)",
          options: [
            {
              name: "--format",
              description: "Output format",
              args: { name: "format", suggestions: ["json", "csv"] },
            },
          ],
        },
      ],
    },

    // ─── teams ───────────────────────────────────────────────────────────
    {
      name: "teams",
      description: "Team management",
      subcommands: [
        {
          name: "list",
          description: "List teams",
          options: [{ name: "--page", args: { name: "n" } }],
        },
        {
          name: "create",
          description: "Create a team",
          options: [
            {
              name: "--name",
              description: "Team name (required)",
              args: { name: "name" },
              isRequired: true,
            },
          ],
        },
        {
          name: "update",
          description: "Update a team (rename)",
          args: idArg,
          options: [{ name: "--name", description: "New team name", args: { name: "name" } }],
        },
        {
          name: "delete",
          description: "Delete a non-imported team (destructive)",
          args: idArg,
          options: [{ name: "--confirm", description: "Required for destructive operation" }],
        },
        {
          name: "link",
          description: "Link a resource to a team",
          args: [
            { name: "team-id" },
            { name: "resource-type", suggestions: resourceTypes },
            { name: "resource-id" },
          ],
        },
        {
          name: "unlink",
          description: "Unlink a resource from a team",
          args: [
            { name: "team-id" },
            { name: "resource-type", suggestions: resourceTypes },
            { name: "resource-id" },
          ],
        },
        {
          name: "remove-user",
          description: "Remove a user from a team (destructive)",
          args: [{ name: "team-id" }, { name: "user-id" }],
          options: [{ name: "--confirm", description: "Required for destructive operation" }],
        },
      ],
    },

    // ─── users ───────────────────────────────────────────────────────────
    {
      name: "users",
      description: "Workspace users",
      subcommands: [
        { name: "list", description: "List users" },
        { name: "get", description: "Get a user", args: idArg },
      ],
    },

    // ─── containers ──────────────────────────────────────────────────────
    {
      name: "containers",
      description: "Container repositories",
      subcommands: [
        { name: "list", description: "List container repositories" },
        { name: "get", description: "Get a container repo", args: idArg },
        {
          name: "sbom",
          description: "Export the SBOM for a container",
          args: idArg,
          options: [
            {
              name: "--format",
              description: "Output format passthrough",
              args: { name: "format" },
            },
          ],
        },
      ],
    },

    // ─── clouds ──────────────────────────────────────────────────────────
    {
      name: "clouds",
      description: "Connected cloud environments",
      subcommands: [
        { name: "list", description: "List connected clouds" },
        { name: "assets", description: "List cloud assets" },
      ],
    },

    // ─── apps / vms / licenses / webhooks ────────────────────────────────
    {
      name: "apps",
      description: "Zen apps",
      subcommands: [{ name: "list", description: "List Zen apps" }],
    },
    {
      name: "vms",
      description: "Virtual machines",
      subcommands: [
        { name: "list", description: "List virtual machines" },
        {
          name: "sbom",
          description: "Export the SBOM for a virtual machine",
          args: idArg,
          options: [
            {
              name: "--format",
              description: "Export format",
              args: { name: "format", suggestions: ["sbom", "sbom_spdx", "csv"] },
            },
          ],
        },
      ],
    },
    {
      name: "licenses",
      description: "License inventory",
      subcommands: [{ name: "list", description: "List licenses across the workspace" }],
    },
    {
      name: "webhooks",
      description: "Configured webhooks",
      subcommands: [{ name: "list", description: "List webhooks" }],
    },

    // ─── activity / pr-checks ────────────────────────────────────────────
    {
      name: "activity",
      description: "Workspace activity log",
      options: [
        { name: "--from", description: "ISO date (inclusive)", args: { name: "date" } },
        { name: "--to", description: "ISO date (inclusive)", args: { name: "date" } },
        { name: "--user", description: "Filter by user type", args: { name: "type" } },
      ],
    },
    {
      name: "pr-checks",
      description: "CI/PR scans",
      subcommands: [
        {
          name: "list",
          description: "List PR checks",
          options: [{ name: "--repo", description: "Filter by repo ID", args: { name: "id" } }],
        },
      ],
    },

    // ─── compliance ──────────────────────────────────────────────────────
    {
      name: "compliance",
      description: "Compliance overviews",
      subcommands: complianceFrameworks.map((f) => ({
        name: f.name as string,
        description: `${(f.name as string).toUpperCase()} compliance overview`,
      })),
    },

    // ─── custom-rules / pentest / tasks ──────────────────────────────────
    {
      name: "custom-rules",
      description: "Custom SAST rules",
      subcommands: [
        { name: "list", description: "List custom rules" },
        { name: "get", description: "Get a custom rule", args: idArg },
      ],
    },
    {
      name: "pentest",
      description: "AI Pentesting",
      subcommands: [
        { name: "get", description: "Get a pentest assessment", args: idArg },
        { name: "attack", description: "Get attack analysis", args: idArg },
      ],
    },
    {
      name: "tasks",
      description: "Task tracker integrations",
      subcommands: [
        { name: "projects", description: "List task tracking projects" },
        {
          name: "list",
          description: "List tasks in a project",
          args: { name: "project-id" },
          options: [{ name: "--search", description: "Search tasks", args: { name: "query" } }],
        },
      ],
    },

    // ─── research + top-level shortcuts ──────────────────────────────────
    {
      name: "research",
      description: "Vulnerability research lookups",
      subcommands: [
        { name: "cve", description: "Get CVE details", args: { name: "cve-id" } },
        {
          name: "changelog",
          description: "Package changelog summary",
          args: { name: "package" },
          options: changelogOptions,
        },
        { name: "malware-packages", description: "Recently flagged malware packages" },
      ],
    },
    { name: "cve", description: "CVE details (shortcut)", args: { name: "cve-id" } },
    {
      name: "changelog",
      description: "Package changelog (shortcut)",
      args: { name: "package" },
      options: changelogOptions,
    },
    { name: "malware-packages", description: "Malware packages (shortcut)" },

    // ─── report ──────────────────────────────────────────────────────────
    {
      name: "report",
      description: "Workspace reports",
      subcommands: [
        {
          name: "pdf",
          description: "Export workspace report as PDF",
          options: [
            {
              name: "--out",
              description: "Write to file instead of stdout",
              args: { name: "path", template: "filepaths" },
            },
            {
              name: "--sections",
              description: "Comma-separated report sections",
              args: { name: "sections" },
              isRequired: true,
            },
            {
              name: "--team",
              description: "Filter report by team ID",
              args: { name: "id" },
            },
          ],
        },
      ],
    },

    // ─── completion (cobra-generated) ────────────────────────────────────
    {
      name: "completion",
      description: "Generate shell completion scripts",
      subcommands: shells.map((s) => ({
        name: s.name as string,
        description: `Generate ${s.name as string} completion script`,
      })),
    },
  ],
  options: globalOptions,
};

export default completion;
