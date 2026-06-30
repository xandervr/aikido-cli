package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

type documentedEndpoint struct {
	Method  string `json:"method" aikido:"column,header=Method"`
	Path    string `json:"path" aikido:"column,header=Path"`
	Summary string `json:"summary" aikido:"column,header=Summary"`
}

var documentedEndpoints = []documentedEndpoint{
	{Method: "POST", Path: "/token", Summary: "Get access token"},
	{Method: "GET", Path: "/open-issue-groups", Summary: "List open issue groups"},
	{Method: "GET", Path: "/issues/export", Summary: "Export all issues"},
	{Method: "GET", Path: "/issues/counts", Summary: "Get issue counts"},
	{Method: "GET", Path: "/issues/detail/bulk", Summary: "Get issue details bulk"},
	{Method: "GET", Path: "/issues/{issue_id}", Summary: "Get issue detail"},
	{Method: "GET", Path: "/issues/{issue_id}/reachability", Summary: "Get issue reachability"},
	{Method: "PUT", Path: "/issues/{issue_id}/snooze", Summary: "Snooze an issue"},
	{Method: "PUT", Path: "/issues/{issue_id}/unsnooze", Summary: "Unsnooze an issue"},
	{Method: "PUT", Path: "/issues/{issue_id}/ignore", Summary: "Ignore an issue"},
	{Method: "PUT", Path: "/issues/{issue_id}/unignore", Summary: "Unignore an issue"},
	{Method: "PUT", Path: "/issues/{issue_id}/solve", Summary: "Solve an issue"},
	{Method: "GET", Path: "/issues/groups/{issue_group_id}", Summary: "Get issue group detail"},
	{Method: "PUT", Path: "/issues/groups/{issue_group_id}/ignore", Summary: "Ignore an issue group"},
	{Method: "PUT", Path: "/issues/groups/{issue_group_id}/unignore", Summary: "Unignore an issue group"},
	{Method: "PUT", Path: "/issues/groups/{issue_group_id}/snooze", Summary: "Snooze an issue group"},
	{Method: "PUT", Path: "/issues/groups/{issue_group_id}/unsnooze", Summary: "Unsnooze an issue group"},
	{Method: "POST", Path: "/issues/groups/{issue_group_id}/notes", Summary: "Add note to issue group"},
	{Method: "GET", Path: "/issues/groups/{issue_group_id}/notes", Summary: "List notes for issue group"},
	{Method: "GET", Path: "/issues/groups/{issue_group_id}/tasks", Summary: "Get issue group tasks"},
	{Method: "POST", Path: "/issues/{issue_id}/severity/adjust", Summary: "Adjust severity of an issue"},
	{Method: "POST", Path: "/issues/groups/{issue_group_id}/severity/adjust", Summary: "Adjust severity of an issue group"},
	{Method: "GET", Path: "/clouds", Summary: "List connected clouds"},
	{Method: "GET", Path: "/clouds/rules", Summary: "List cloud rules"},
	{Method: "POST", Path: "/clouds/assets", Summary: "List cloud assets"},
	{Method: "DELETE", Path: "/clouds/{cloud_id}", Summary: "Remove cloud"},
	{Method: "POST", Path: "/clouds/aws", Summary: "Connect AWS cloud"},
	{Method: "PUT", Path: "/clouds/azure/{cloud_id}/credentials", Summary: "Update Azure cloud"},
	{Method: "POST", Path: "/clouds/azure", Summary: "Connect Azure cloud"},
	{Method: "POST", Path: "/clouds/gcp", Summary: "Connect GCP cloud"},
	{Method: "POST", Path: "/clouds/kubernetes", Summary: "Create Kubernetes cloud"},
	{Method: "GET", Path: "/workspace", Summary: "Get workspace info"},
	{Method: "GET", Path: "/workspace/configurationErrors", Summary: "Get workspace configuration errors"},
	{Method: "GET", Path: "/workspace/slaSettings", Summary: "Get SLA settings"},
	{Method: "GET", Path: "/localscan/latest", Summary: "Get latest local scanner version"},
	{Method: "GET", Path: "/endpoint-protection/activityLogs", Summary: "List endpoint activity logs"},
	{Method: "GET", Path: "/endpoint-protection/devices", Summary: "List endpoint devices"},
	{Method: "GET", Path: "/endpoint-protection/installed-packages", Summary: "List installed packages"},
	{Method: "GET", Path: "/endpoint-protection/permission-groups", Summary: "List endpoint permission groups"},
	{Method: "GET", Path: "/endpoint-protection/{ecosystem}/exceptions", Summary: "List endpoint exceptions"},
	{Method: "POST", Path: "/endpoint-protection/{ecosystem}/exceptions", Summary: "Add an endpoint exception"},
	{Method: "DELETE", Path: "/endpoint-protection/exceptions/{package_exception_id}", Summary: "Remove an endpoint exception"},
	{Method: "GET", Path: "/report/export/pdf", Summary: "Export PDF report"},
	{Method: "GET", Path: "/report/activityLog", Summary: "List activity log"},
	{Method: "GET", Path: "/report/ciScans", Summary: "List PR Checks"},
	{Method: "GET", Path: "/report/ciScans/issueActions", Summary: "List PR Check Manual Actions"},
	{Method: "GET", Path: "/report/soc2/overview", Summary: "SOC2 compliance"},
	{Method: "GET", Path: "/report/nis2/overview", Summary: "NIS2 compliance"},
	{Method: "GET", Path: "/report/iso/overview", Summary: "ISO 27001 compliance"},
	{Method: "GET", Path: "/report/cis/overview", Summary: "CIS compliance"},
	{Method: "GET", Path: "/report/cis_aws/overview", Summary: "CIS AWS compliance"},
	{Method: "GET", Path: "/research/malware/packages", Summary: "Get malware packages"},
	{Method: "GET", Path: "/repositories/code", Summary: "List code repositories"},
	{Method: "GET", Path: "/repositories/code/{code_repo_id}", Summary: "Get code repository detail"},
	{Method: "DELETE", Path: "/repositories/code/{code_repo_id}", Summary: "Delete code repo"},
	{Method: "POST", Path: "/repositories/code/{code_repo_id}/scan", Summary: "Scan code repo"},
	{Method: "GET", Path: "/repositories/code/{code_repo_id}/licenses/export", Summary: "Export SBOM"},
	{Method: "PUT", Path: "/repositories/code/{code_repo_id}/devdep-scan", Summary: "Manage dev dep scanning"},
	{Method: "PUT", Path: "/repositories/code/{code_repo_id}/sensitivity", Summary: "Update sensitivity"},
	{Method: "PUT", Path: "/repositories/code/{code_repo_id}/connectivity", Summary: "Update connectivity"},
	{Method: "POST", Path: "/repositories/code/{code_repo_id}/exclude-path", Summary: "Add an exclude path to a code repo"},
	{Method: "POST", Path: "/repositories/code/{code_repo_id}/exclude-path/remove", Summary: "Remove an exclude path from a code repo"},
	{Method: "POST", Path: "/repositories/code/{code_repo_id}/labels", Summary: "Add code repo label"},
	{Method: "POST", Path: "/repositories/code/{code_repo_id}/labels/{label_id}", Summary: "Update code repo label"},
	{Method: "DELETE", Path: "/repositories/code/{code_repo_id}/labels/{label_id}", Summary: "Remove code repo label"},
	{Method: "GET", Path: "/repositories/code/team/{team_id}/licenses/export", Summary: "Export SBOM For Team"},
	{Method: "POST", Path: "/repositories/code/activate", Summary: "Activate code repo"},
	{Method: "POST", Path: "/repositories/code/deactivate", Summary: "Deactivate code repo"},
	{Method: "POST", Path: "/repositories/code/clone", Summary: "Clone code repo"},
	{Method: "POST", Path: "/repositories/code/private-registries", Summary: "Manage private registry"},
	{Method: "POST", Path: "/repositories/import", Summary: "Trigger repositories sync"},
	{Method: "GET", Path: "/repositories/code/sast/rules", Summary: "List SAST rules"},
	{Method: "GET", Path: "/repositories/code/iac/rules", Summary: "List IaC rules"},
	{Method: "GET", Path: "/repositories/code/mobile/rules", Summary: "List Mobile rules"},
	{Method: "GET", Path: "/repositories/sast/custom-rules", Summary: "List custom rules"},
	{Method: "POST", Path: "/repositories/sast/custom-rules", Summary: "Create custom rule"},
	{Method: "GET", Path: "/repositories/sast/custom-rules/{rule_id}", Summary: "Get a custom rule"},
	{Method: "PUT", Path: "/repositories/sast/custom-rules/{rule_id}", Summary: "Edit custom rule"},
	{Method: "DELETE", Path: "/repositories/sast/custom-rules/{rule_id}", Summary: "Remove custom rule"},
	{Method: "POST", Path: "/repositories/code/continuous_integration/checks", Summary: "Configure PR Checks"},
	{Method: "GET", Path: "/containers", Summary: "List containers"},
	{Method: "GET", Path: "/containers/{container_repo_id}", Summary: "Get container"},
	{Method: "GET", Path: "/containers/{container_repo_id}/runners", Summary: "List container runners"},
	{Method: "DELETE", Path: "/containers/{container_repo_id}", Summary: "Delete container"},
	{Method: "GET", Path: "/containers/{container_repo_id}/licenses/export", Summary: "Export SBOM"},
	{Method: "GET", Path: "/containers/{container_repo_id}/sbom/exportRaw", Summary: "Export Raw SBOM"},
	{Method: "PUT", Path: "/containers/{container_repo_id}/sensitivity", Summary: "Update sensitivity"},
	{Method: "PUT", Path: "/containers/{container_repo_id}/internetConnection", Summary: "Update connectivity"},
	{Method: "POST", Path: "/containers/sbom", Summary: "Upload container SBOM"},
	{Method: "POST", Path: "/containers/sbom/generate", Summary: "Generate bulk SBOM"},
	{Method: "POST", Path: "/containers/activate", Summary: "Activate container"},
	{Method: "POST", Path: "/containers/deactivate", Summary: "Deactivate container"},
	{Method: "POST", Path: "/containers/linkCodeRepo", Summary: "Link code repository to container"},
	{Method: "POST", Path: "/containers/unlinkCodeRepo", Summary: "Unlink code repository from container"},
	{Method: "POST", Path: "/containers/updateTagFilter", Summary: "Update container tag filter"},
	{Method: "POST", Path: "/containers/public", Summary: "Add public container"},
	{Method: "POST", Path: "/containers/clone", Summary: "Clone container"},
	{Method: "POST", Path: "/containers/{container_repo_id}/scan", Summary: "Scan container"},
	{Method: "GET", Path: "/containers/registries/{registry_id}", Summary: "Get container registry"},
	{Method: "POST", Path: "/containers/registries/acr", Summary: "Add Azure container registry"},
	{Method: "POST", Path: "/containers/registries/gcp-artifact-registry", Summary: "Add GCP Artifact Registry"},
	{Method: "POST", Path: "/domains", Summary: "Create domain"},
	{Method: "GET", Path: "/domains", Summary: "List domains"},
	{Method: "DELETE", Path: "/domains/{domain_id}", Summary: "Remove domain"},
	{Method: "POST", Path: "/domains/{domain_id}/headers", Summary: "Update Auth Headers"},
	{Method: "POST", Path: "/domains/{domain_id}/custom-headers", Summary: "Update Custom Scan Headers"},
	{Method: "GET", Path: "/domains/{domain_id}/subdomains", Summary: "List Subdomains"},
	{Method: "POST", Path: "/domains/{domain_id}/subdomains", Summary: "Add Subdomain"},
	{Method: "PUT", Path: "/domains/{domain_id}/update/openapi-spec", Summary: "Update OpenAPI spec"},
	{Method: "POST", Path: "/domains/scan", Summary: "Start scan for a domain"},
	{Method: "GET", Path: "/teams", Summary: "List teams"},
	{Method: "POST", Path: "/teams", Summary: "Create team"},
	{Method: "PUT", Path: "/teams/{team_id}", Summary: "Update team"},
	{Method: "DELETE", Path: "/teams/{team_id}", Summary: "Delete team"},
	{Method: "POST", Path: "/teams/{team_id}/linkResource", Summary: "Link resource to team"},
	{Method: "POST", Path: "/teams/{team_id}/unlinkResource", Summary: "Unlink resource from team"},
	{Method: "POST", Path: "/teams/{team_id}/addUser", Summary: "Add user to team"},
	{Method: "POST", Path: "/teams/{team_id}/removeUser", Summary: "Remove user from team"},
	{Method: "GET", Path: "/users", Summary: "List users"},
	{Method: "GET", Path: "/users/{user_id}", Summary: "Get user"},
	{Method: "GET", Path: "/users/ide/adoption", Summary: "List IDE adoption"},
	{Method: "PUT", Path: "/users/{user_id}/rights", Summary: "Update user rights"},
	{Method: "GET", Path: "/firewall/apps", Summary: "List apps"},
	{Method: "POST", Path: "/firewall/apps", Summary: "Create app"},
	{Method: "GET", Path: "/firewall/apps/{app_id}", Summary: "Get app"},
	{Method: "PUT", Path: "/firewall/apps/{app_id}", Summary: "Update app"},
	{Method: "DELETE", Path: "/firewall/apps/{app_id}", Summary: "Delete app"},
	{Method: "POST", Path: "/firewall/apps/{app_id}/token", Summary: "Rotate app token"},
	{Method: "PUT", Path: "/firewall/apps/{service_id}/blocking", Summary: "Update blocking mode"},
	{Method: "PUT", Path: "/firewall/{app_id}/users/{user_id}", Summary: "Update user"},
	{Method: "GET", Path: "/firewall/apps/{app_id}/users", Summary: "List users"},
	{Method: "PUT", Path: "/firewall/apps/{app_id}/ip-blocklist", Summary: "Update IP blocklist"},
	{Method: "GET", Path: "/firewall/apps/{app_id}/bot-lists", Summary: "Get bot lists"},
	{Method: "PUT", Path: "/firewall/apps/{app_id}/bot-lists", Summary: "Update bot lists"},
	{Method: "GET", Path: "/firewall/apps/{app_id}/ip-lists", Summary: "Get threat lists"},
	{Method: "PUT", Path: "/firewall/apps/{app_id}/ip-lists", Summary: "Update threat lists"},
	{Method: "GET", Path: "/firewall/apps/{app_id}/countries", Summary: "Get countries"},
	{Method: "PUT", Path: "/firewall/apps/{app_id}/countries", Summary: "Update countries"},
	{Method: "GET", Path: "/firewall/apps/{app_id}/events/{event_id}", Summary: "Get event"},
	{Method: "GET", Path: "/virtual-machines", Summary: "List virtual machines"},
	{Method: "GET", Path: "/virtual-machines/{virtual_machine_id}/export/{format}", Summary: "Export SBOM"},
	{Method: "GET", Path: "/task_tracking/projects", Summary: "List task tracking projects"},
	{Method: "GET", Path: "/task_tracking/projects/{project_id}/tasks", Summary: "List tasks from project"},
	{Method: "GET", Path: "/task_tracking/integrations", Summary: "List task tracking integrations"},
	{Method: "POST", Path: "/task_tracking/mapCodeReposToProjects", Summary: "Map code repo to task tracking projects"},
	{Method: "POST", Path: "/task_tracking/linkTaskToIssueGroup", Summary: "Link existing task to issue"},
	{Method: "GET", Path: "/task_tracking/projectMapping", Summary: "Get project mapping"},
	{Method: "GET", Path: "/pentests/assessments/{assessment_id}/detail", Summary: "Get pentest assessment"},
	{Method: "POST", Path: "/pentests/assessments/createDraft", Summary: "Create pentest draft"},
	{Method: "GET", Path: "/pentests/issues/{issue_id}/attackAnalysis", Summary: "Get attack analysis"},
	{Method: "GET", Path: "/licenses", Summary: "List & Search SBOM"},
	{Method: "POST", Path: "/licenses/overwrite", Summary: "Overwrite License"},
	{Method: "GET", Path: "/changelog-summary", Summary: "Get changelog summary"},
	{Method: "GET", Path: "/cve/{cve_id}", Summary: "Get CVE details"},
	{Method: "POST", Path: "/bug_bounty/program/{program_id}/report", Summary: "Validate Bug Bounty Report"},
	{Method: "POST", Path: "/access-tokens/code-scanning", Summary: "Update Code Scanning Access Token"},
	{Method: "GET", Path: "/openapi/spec", Summary: "Get OpenAPI spec"},
	{Method: "GET", Path: "/code-quality/findings", Summary: "List code quality findings for a pull request"},
	{Method: "GET", Path: "/webhooks", Summary: "List webhooks"},
	{Method: "POST", Path: "/webhooks", Summary: "Add webhook"},
	{Method: "DELETE", Path: "/webhooks/{webhook_id}", Summary: "Remove webhook"},
}

type apiRequestOptions struct {
	query    []string
	body     string
	bodyFile string
	out      string
}

type endpointCommandConfig struct {
	Use     string
	Short   string
	Method  string
	Args    cobra.PositionalArgs
	Path    func([]string) string
	Confirm bool
}

func NewAPI(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "api", Short: "Call any documented Aikido public API endpoint"}
	cmd.AddCommand(apiEndpoints(g))
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		cmd.AddCommand(apiMethod(g, method))
	}
	return cmd
}

func apiEndpoints(g *cli.Globals) *cobra.Command {
	var method, search string
	cmd := &cobra.Command{
		Use:   "endpoints",
		Short: "List the current documented endpoint catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			filtered := make([]documentedEndpoint, 0, len(documentedEndpoints))
			method = strings.ToUpper(method)
			search = strings.ToLower(search)
			for _, ep := range documentedEndpoints {
				if method != "" && ep.Method != method {
					continue
				}
				if search != "" {
					haystack := strings.ToLower(ep.Path + " " + ep.Summary)
					if !strings.Contains(haystack, search) {
						continue
					}
				}
				filtered = append(filtered, ep)
			}
			return g.Renderer().Render(filtered)
		},
	}
	cmd.Flags().StringVar(&method, "method", "", "filter by HTTP method")
	cmd.Flags().StringVar(&search, "search", "", "filter by path or summary text")
	return cmd
}

func apiMethod(g *cli.Globals, method string) *cobra.Command {
	var opts apiRequestOptions
	cmd := &cobra.Command{
		Use:   strings.ToLower(method) + " <path>",
		Short: method + " an arbitrary Aikido public API path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPIRequest(cmd.Context(), g, method, args[0], opts)
		},
	}
	addAPIRequestFlags(cmd, &opts, method != http.MethodGet && method != http.MethodDelete)
	return cmd
}

func endpointCommand(g *cli.Globals, cfg endpointCommandConfig) *cobra.Command {
	var opts apiRequestOptions
	var confirm bool
	if cfg.Args == nil {
		cfg.Args = cobra.NoArgs
	}
	cmd := &cobra.Command{
		Use:   cfg.Use,
		Short: cfg.Short,
		Args:  cfg.Args,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Confirm && !confirm {
				return &cli.ExitError{Code: cli.ExitUsage, Err: errors.New("destructive: pass --confirm")}
			}
			return runAPIRequest(cmd.Context(), g, cfg.Method, cfg.Path(args), opts)
		},
	}
	addAPIRequestFlags(cmd, &opts, cfg.Method != http.MethodGet && cfg.Method != http.MethodDelete)
	if cfg.Confirm {
		cmd.Flags().BoolVar(&confirm, "confirm", false, "required for destructive operation")
	}
	return cmd
}

func addAPIRequestFlags(cmd *cobra.Command, opts *apiRequestOptions, includeBody bool) {
	cmd.Flags().StringArrayVarP(&opts.query, "query", "q", nil, "query parameter as key=value; repeatable")
	if includeBody {
		cmd.Flags().StringVar(&opts.body, "body", "", "JSON request body")
		cmd.Flags().StringVar(&opts.bodyFile, "body-file", "", "read JSON request body from file")
	}
	cmd.Flags().StringVar(&opts.out, "out", "", "write raw response bytes to file")
}

func runAPIRequest(ctx context.Context, g *cli.Globals, method, path string, opts apiRequestOptions) error {
	query, err := parseQueryPairs(opts.query)
	if err != nil {
		return &cli.ExitError{Code: cli.ExitUsage, Err: err}
	}
	body, err := parseJSONBody(opts.body, opts.bodyFile)
	if err != nil {
		return &cli.ExitError{Code: cli.ExitUsage, Err: err}
	}
	c, err := g.Client()
	if err != nil {
		return err
	}
	resp, _, err := c.Raw(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	if opts.out != "" {
		return os.WriteFile(opts.out, resp, 0o644)
	}
	_, err = fmt.Fprint(g.Renderer().Out, string(resp))
	return err
}

func parseQueryPairs(pairs []string) (map[string]string, error) {
	query := map[string]string{}
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("--query must be key=value, got %q", pair)
		}
		query[key] = value
	}
	if len(query) == 0 {
		return nil, nil
	}
	return query, nil
}

func parseJSONBody(body, bodyFile string) (any, error) {
	if body != "" && bodyFile != "" {
		return nil, errors.New("use only one of --body or --body-file")
	}
	if bodyFile != "" {
		b, err := os.ReadFile(bodyFile)
		if err != nil {
			return nil, err
		}
		body = string(b)
	}
	if body == "" {
		return nil, nil
	}
	var raw any
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		return nil, fmt.Errorf("body must be valid JSON: %w", err)
	}
	return raw, nil
}

func staticPath(path string) func([]string) string {
	return func([]string) string {
		return path
	}
}

func oneArgPath(format string) func([]string) string {
	return func(args []string) string {
		return fmt.Sprintf(format, args[0])
	}
}
