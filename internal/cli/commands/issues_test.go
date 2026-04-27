package commands

import (
	"encoding/json"
	"testing"
)

func TestIssueGroup_UnmarshalProbesAliases(t *testing.T) {
	cases := []struct {
		name string
		in   string
		repo string
		stat string
	}{
		{"flat repo_name + status",
			`{"id":1,"title":"x","severity":"high","type":"sast","repo_name":"acme","status":"open"}`,
			"acme", "open"},
		{"nested code_repo.name + state",
			`{"id":2,"title":"y","severity":"low","issue_type":"iac","code_repo":{"name":"infra"},"state":"snoozed"}`,
			"infra", "snoozed"},
		{"is_open boolean",
			`{"id":3,"title":"z","severity":"medium","is_open":true}`,
			"", "open"},
		{"group_status (Aikido actual key)",
			`{"id":5,"title":"a","severity":"high","type":"open_source","group_status":"open"}`,
			"", "open"},
		{"repo from locations[0].code_repo_name",
			`{"id":6,"title":"b","severity":"low","type":"sast","group_status":"open","locations":[{"code_repo_name":"acme-api"}]}`,
			"acme-api", "open"},
		{"repo from locations counts distinct repos only",
			`{"id":7,"title":"c","severity":"low","type":"open_source","group_status":"open","locations":[{"code_repo_name":"a"},{"code_repo_name":"b"},{"code_repo_name":"c"}]}`,
			"a (+2)", "open"},
		{"locations with duplicate repos dedupe to single repo",
			`{"id":8,"title":"d","severity":"low","type":"open_source","group_status":"open","locations":[{"code_repo_name":"a"},{"code_repo_name":"a"},{"code_repo_name":"a"}]}`,
			"a", "open"},
		{"locations with mixed dupes count only distinct repos",
			`{"id":9,"title":"e","severity":"low","type":"open_source","group_status":"open","locations":[{"code_repo_name":"a"},{"code_repo_name":"a"},{"code_repo_name":"b"},{"code_repo_name":"a"}]}`,
			"a (+1)", "open"},
		{"repo from nested locations[0].repository.name",
			`{"id":8,"title":"d","severity":"low","type":"open_source","group_status":"snoozed","locations":[{"repository":{"name":"infra"}}]}`,
			"infra", "snoozed"},
		{"ignored boolean",
			`{"id":4,"title":"w","severity":"low","ignored":true}`,
			"", "ignored"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var g IssueGroup
			if err := json.Unmarshal([]byte(tc.in), &g); err != nil {
				t.Fatal(err)
			}
			if g.Repo != tc.repo {
				t.Errorf("repo = %q, want %q", g.Repo, tc.repo)
			}
			if g.Status != tc.stat {
				t.Errorf("status = %q, want %q", g.Status, tc.stat)
			}
		})
	}
}

func TestIssueGroup_MarshalPreservesRaw(t *testing.T) {
	in := `{"id":42,"title":"x","random_field":"keep me","status":"open"}`
	var g IssueGroup
	if err := json.Unmarshal([]byte(in), &g); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["random_field"] != "keep me" {
		t.Fatalf("lost random_field; got %+v", got)
	}
}
