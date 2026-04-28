package commands

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

type ctor func(*cli.Globals) *cobra.Command

func TestSimpleListCommands_HitExpectedPath(t *testing.T) {
	type commandCase struct {
		name       string
		factory    ctor
		args       []string
		wantPath   string
		wantMethod string
		wantQuery  string
	}
	cmdCase := func(name string, factory ctor, args []string, wantPath, wantMethod string) commandCase {
		return commandCase{name: name, factory: factory, args: args, wantPath: wantPath, wantMethod: wantMethod}
	}
	cmdCaseWithQuery := func(name string, factory ctor, args []string, wantPath, wantMethod, wantQuery string) commandCase {
		return commandCase{name: name, factory: factory, args: args, wantPath: wantPath, wantMethod: wantMethod, wantQuery: wantQuery}
	}

	cases := []commandCase{
		cmdCase("workspace info", NewWorkspace, []string{"info"}, "/workspace", "GET"),
		cmdCase("workspace config-errors", NewWorkspace, []string{"config-errors"}, "/workspace/configuration-errors", "GET"),
		cmdCase("workspace introspect", NewWorkspace, []string{"introspect"}, "/openapi/spec", "GET"),
		cmdCase("users list", NewUsers, []string{"list"}, "/users", "GET"),
		cmdCase("users get", NewUsers, []string{"get", "42"}, "/users/42", "GET"),
		cmdCase("containers list", NewContainers, []string{"list"}, "/containers", "GET"),
		cmdCase("containers get", NewContainers, []string{"get", "9"}, "/containers/9", "GET"),
		cmdCase("containers sbom", NewContainers, []string{"sbom", "9"}, "/containers/9/licenses/export", "GET"),
		cmdCase("clouds list", NewClouds, []string{"list"}, "/clouds", "GET"),
		cmdCase("clouds assets", NewClouds, []string{"assets"}, "/clouds/assets", "POST"),
		cmdCase("apps list", NewApps, []string{"list"}, "/firewall/apps", "GET"),
		cmdCase("vms list", NewVMs, []string{"list"}, "/virtual-machines", "GET"),
		cmdCase("vms sbom", NewVMs, []string{"sbom", "9"}, "/virtual-machines/9/export/sbom", "GET"),
		cmdCase("licenses list", NewLicenses, []string{"list"}, "/licenses", "GET"),
		cmdCase("webhooks list", NewWebhooks, []string{"list"}, "/webhooks", "GET"),
		cmdCase("pr-checks list", NewPRChecks, []string{"list"}, "/report/ciScans", "GET"),
		cmdCaseWithQuery("pr-checks list repo", NewPRChecks, []string{"list", "--repo", "12"}, "/report/ciScans", "GET", "filter_code_repo_id=12"),
		cmdCase("compliance soc2", NewCompliance, []string{"soc2"}, "/report/soc2/overview", "GET"),
		cmdCase("compliance nis2", NewCompliance, []string{"nis2"}, "/report/nis2/overview", "GET"),
		cmdCase("compliance iso27001", NewCompliance, []string{"iso27001"}, "/report/iso/overview", "GET"),
		cmdCase("custom-rules list", NewCustomRules, []string{"list"}, "/repositories/sast/custom-rules", "GET"),
		cmdCase("custom-rules get", NewCustomRules, []string{"get", "5"}, "/repositories/sast/custom-rules/5", "GET"),
		cmdCase("pentest get", NewPentest, []string{"get", "550e8400-e29b-41d4-a716-446655440000"}, "/pentests/assessments/550e8400-e29b-41d4-a716-446655440000/detail", "GET"),
		cmdCase("pentest attack", NewPentest, []string{"attack", "2"}, "/pentests/issues/2/attackAnalysis", "GET"),
		cmdCase("tasks projects", NewTasks, []string{"projects"}, "/task_tracking/projects", "GET"),
		cmdCase("tasks list", NewTasks, []string{"list", "7"}, "/task_tracking/projects/7/tasks", "GET"),
		cmdCase("research cve", NewResearch, []string{"cve", "CVE-2026-1234"}, "/cve/CVE-2026-1234", "GET"),
		cmdCaseWithQuery("research changelog", NewResearch, []string{"changelog", "jsonpath-plus", "--from", "5.1.0", "--to", "10.2.0", "--language", "JS"}, "/changelog-summary", "GET", "from_version=5.1.0&language=JS&package_name=jsonpath-plus&to_version=10.2.0"),
		cmdCase("research malware-packages", NewResearch, []string{"malware-packages"}, "/research/malware/packages", "GET"),
		cmdCase("cve shortcut", NewCVE, []string{"CVE-2026-1234"}, "/cve/CVE-2026-1234", "GET"),
		cmdCaseWithQuery("changelog shortcut", NewChangelog, []string{"jsonpath-plus", "--from", "5.1.0", "--to", "10.2.0", "--language", "JS"}, "/changelog-summary", "GET", "from_version=5.1.0&language=JS&package_name=jsonpath-plus&to_version=10.2.0"),
		cmdCase("malware-packages shortcut", NewMalwarePackages, []string{}, "/research/malware/packages", "GET"),
		cmdCase("repos list", NewRepos, []string{"list"}, "/repositories/code", "GET"),
		cmdCaseWithQuery("repos list search", NewRepos, []string{"list", "--search", "api"}, "/repositories/code", "GET", "filter_name=api"),
		cmdCase("repos get", NewRepos, []string{"get", "3"}, "/repositories/code/3", "GET"),
		cmdCase("repos sbom", NewRepos, []string{"sbom", "3"}, "/repositories/code/3/licenses/export", "GET"),
		cmdCase("issues list", NewIssues, []string{"list"}, "/open-issue-groups", "GET"),
		cmdCase("issues get", NewIssues, []string{"get", "11"}, "/issues/groups/11", "GET"),
		cmdCase("teams list", NewTeams, []string{"list"}, "/teams", "GET"),
		cmdCase("teams link", NewTeams, []string{"link", "4", "repo", "9"}, "/teams/4/linkResource", "POST"),
		cmdCase("teams unlink", NewTeams, []string{"unlink", "4", "repo", "9"}, "/teams/4/unlinkResource", "POST"),
		cmdCase("teams remove-user", NewTeams, []string{"remove-user", "4", "7", "--confirm"}, "/teams/4/removeUser", "POST"),
		cmdCase("activity top-level", NewActivity, []string{}, "/report/activityLog", "GET"),
		cmdCaseWithQuery("activity dates", NewActivity, []string{"--from", "2026-01-01", "--to", "2026-01-31"}, "/report/activityLog", "GET", "end=2026-01-31&start=2026-01-01"),
		cmdCaseWithQuery("report pdf", NewReport, []string{"pdf", "--sections", "soc2"}, "/report/export/pdf", "GET", "included_sections=soc2"),
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath, gotMethod, gotQuery string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotMethod = r.Method
				gotQuery = r.URL.RawQuery
				if r.Method == "GET" {
					_, _ = w.Write([]byte(`[]`))
				} else {
					_, _ = w.Write([]byte(`{}`))
				}
			}))
			defer srv.Close()
			t.Setenv("AIKIDO_ACCESS_TOKEN", "test")
			t.Setenv("AIKIDO_BASE_URL", srv.URL)

			root, g := cli.NewRoot()
			child := tc.factory(g)
			root.AddCommand(child)
			args := append([]string{child.Name()}, tc.args...)
			args = append(args, "--json")
			root.SetArgs(args)
			if err := root.Execute(); err != nil {
				t.Fatalf("execute %v: %v", args, err)
			}
			if gotPath != tc.wantPath {
				t.Errorf("path = %q, want %q", gotPath, tc.wantPath)
			}
			if gotMethod != tc.wantMethod {
				t.Errorf("method = %q, want %q", gotMethod, tc.wantMethod)
			}
			if gotQuery != tc.wantQuery {
				t.Errorf("query = %q, want %q", gotQuery, tc.wantQuery)
			}
		})
	}
}

func TestTeamsCreate_PostsName(t *testing.T) {
	var gotMethod, gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		_, _ = w.Write([]byte(`{"id":1,"name":"Platform"}`))
	}))
	defer srv.Close()
	t.Setenv("AIKIDO_ACCESS_TOKEN", "k")
	t.Setenv("AIKIDO_BASE_URL", srv.URL)

	root, g := cli.NewRoot()
	root.AddCommand(NewTeams(g))
	root.SetArgs([]string{"teams", "create", "--name", "Platform", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if gotMethod != "POST" || gotPath != "/teams" {
		t.Fatalf("expected POST /teams, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"name":"Platform"`) {
		t.Fatalf("body missing name field: %q", gotBody)
	}
}

func TestTeamsResourceCommands_PostDocumentedBody(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantPath  string
		wantField string
	}{
		{"link repo", []string{"teams", "link", "4", "repo", "9", "--json"}, "/teams/4/linkResource", `"repo_id":9`},
		{"unlink app", []string{"teams", "unlink", "4", "app", "8", "--json"}, "/teams/4/unlinkResource", `"zen_app_id":8`},
		{"remove user", []string{"teams", "remove-user", "4", "7", "--confirm", "--json"}, "/teams/4/removeUser", `"user_id":7`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotMethod, gotPath, gotBody string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.Path
				body, _ := io.ReadAll(r.Body)
				gotBody = string(body)
				_, _ = w.Write([]byte(`{"success":1}`))
			}))
			defer srv.Close()
			t.Setenv("AIKIDO_ACCESS_TOKEN", "k")
			t.Setenv("AIKIDO_BASE_URL", srv.URL)

			root, g := cli.NewRoot()
			root.AddCommand(NewTeams(g))
			root.SetArgs(tc.args)
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if gotMethod != "POST" || gotPath != tc.wantPath {
				t.Fatalf("expected POST %s, got %s %s", tc.wantPath, gotMethod, gotPath)
			}
			if !strings.Contains(gotBody, tc.wantField) {
				t.Fatalf("body missing %s: %q", tc.wantField, gotBody)
			}
		})
	}
}
