package commands

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTeamJSONDoesNotInventUserCount(t *testing.T) {
	raw := []byte(`[
		{
			"id": 1,
			"name": "Platform",
			"external_source": null,
			"external_source_id": null,
			"responsibilities": [{"id": 5, "type": "code_repository"}],
			"active": true
		}
	]`)

	var teams []Team
	if err := json.Unmarshal(raw, &teams); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(teams)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "user_count") {
		t.Fatalf("team JSON invented user_count: %s", got)
	}

	var decoded []map[string]any
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("invalid json: %v\nbody=%s", err, got)
	}
	if _, ok := decoded[0]["active"]; !ok {
		t.Fatalf("team JSON dropped raw API fields: %s", got)
	}
}
