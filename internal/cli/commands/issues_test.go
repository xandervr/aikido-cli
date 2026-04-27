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
