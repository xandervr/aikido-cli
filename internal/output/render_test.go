package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type sample struct {
	ID   int    `json:"id"   aikido:"column,header=ID"`
	Name string `json:"name" aikido:"column,header=Name"`
}

func TestRenderJSON_PipesArray(t *testing.T) {
	var buf bytes.Buffer
	r := &Renderer{Out: &buf, ForceJSON: true}
	if err := r.Render([]sample{{1, "alpha"}, {2, "beta"}}); err != nil {
		t.Fatal(err)
	}
	var got []sample
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v\nbody=%s", err, buf.String())
	}
	if len(got) != 2 || got[0].Name != "alpha" {
		t.Fatalf("bad payload: %+v", got)
	}
}

func TestRenderTable_HeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	r := &Renderer{Out: &buf, ForceTable: true}
	if err := r.Render([]sample{{1, "alpha"}, {2, "beta"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "Name") {
		t.Fatalf("missing headers: %q", out)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Fatalf("missing rows: %q", out)
	}
}

func TestRenderTable_FallsBackToJSONWhenNoTags(t *testing.T) {
	type unTagged struct{ A, B string }
	var buf bytes.Buffer
	r := &Renderer{Out: &buf, ForceTable: true}
	if err := r.Render(unTagged{A: "x", B: "y"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"A": "x"`) {
		t.Fatalf("expected JSON fallback, got: %q", buf.String())
	}
}
