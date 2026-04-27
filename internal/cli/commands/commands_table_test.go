package commands

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

type ctor func(*cli.Globals) *cobra.Command

func TestSimpleListCommands_HitExpectedPath(t *testing.T) {
	cases := []struct {
		name       string
		factory    ctor
		args       []string
		wantPath   string
		wantMethod string
	}{
		{"workspace info", NewWorkspace, []string{"info"}, "/workspace", "GET"},
		{"workspace config-errors", NewWorkspace, []string{"config-errors"}, "/workspace/configuration-errors", "GET"},
		{"users list", NewUsers, []string{"list"}, "/users", "GET"},
		{"users get", NewUsers, []string{"get", "42"}, "/users/42", "GET"},
		{"containers list", NewContainers, []string{"list"}, "/repositories/container", "GET"},
		{"containers get", NewContainers, []string{"get", "9"}, "/repositories/container/9", "GET"},
		{"clouds list", NewClouds, []string{"list"}, "/clouds", "GET"},
		{"clouds assets", NewClouds, []string{"assets"}, "/clouds/assets", "POST"},
		{"apps list", NewApps, []string{"list"}, "/apps", "GET"},
		{"vms list", NewVMs, []string{"list"}, "/virtual-machines", "GET"},
		{"licenses list", NewLicenses, []string{"list"}, "/licenses", "GET"},
		{"webhooks list", NewWebhooks, []string{"list"}, "/webhooks", "GET"},
		{"pr-checks list", NewPRChecks, []string{"list"}, "/ci-scans", "GET"},
		{"compliance soc2", NewCompliance, []string{"soc2"}, "/compliance/soc2", "GET"},
		{"compliance nis2", NewCompliance, []string{"nis2"}, "/compliance/nis2", "GET"},
		{"compliance iso27001", NewCompliance, []string{"iso27001"}, "/compliance/iso27001", "GET"},
		{"custom-rules list", NewCustomRules, []string{"list"}, "/custom-rules", "GET"},
		{"custom-rules get", NewCustomRules, []string{"get", "5"}, "/custom-rules/5", "GET"},
		{"pentest get", NewPentest, []string{"get", "1"}, "/pentests/1", "GET"},
		{"pentest attack", NewPentest, []string{"attack", "2"}, "/pentests/attacks/2", "GET"},
		{"tasks projects", NewTasks, []string{"projects"}, "/tasks/projects", "GET"},
		{"tasks list", NewTasks, []string{"list", "7"}, "/tasks/projects/7", "GET"},
		{"research malware-packages", NewResearch, []string{"malware-packages"}, "/research/malware-packages", "GET"},
		{"repos list", NewRepos, []string{"list"}, "/repositories/code", "GET"},
		{"repos get", NewRepos, []string{"get", "3"}, "/repositories/code/3", "GET"},
		{"issues list", NewIssues, []string{"list"}, "/open-issue-groups", "GET"},
		{"issues get", NewIssues, []string{"get", "11"}, "/open-issue-groups/11", "GET"},
		{"teams list", NewTeams, []string{"list"}, "/teams", "GET"},
		{"activity top-level", NewActivity, []string{}, "/activity-log", "GET"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath, gotMethod string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotMethod = r.Method
				if r.Method == "GET" {
					_, _ = w.Write([]byte(`[]`))
				} else {
					_, _ = w.Write([]byte(`{}`))
				}
			}))
			defer srv.Close()
			t.Setenv("AIKIDO_API_KEY", "test")
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
	t.Setenv("AIKIDO_API_KEY", "k")
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
