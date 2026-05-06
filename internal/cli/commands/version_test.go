package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/version"
)

func withVersionInfo(t *testing.T, v, commit, date string) {
	t.Helper()
	oldVersion, oldCommit, oldDate := version.Version, version.Commit, version.Date
	version.Version = v
	version.Commit = commit
	version.Date = date
	t.Cleanup(func() {
		version.Version = oldVersion
		version.Commit = oldCommit
		version.Date = oldDate
	})
}

func TestVersionCommandPrintsPlainVersion(t *testing.T) {
	withVersionInfo(t, "1.2.3", "abc123", "2026-05-06T08:00:00Z")

	root, g := cli.NewRoot()
	root.AddCommand(NewVersion(g))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != "aikido version 1.2.3" {
		t.Fatalf("version output = %q", got)
	}
}

func TestVersionCommandPrintsJSON(t *testing.T) {
	withVersionInfo(t, "1.2.3", "abc123", "2026-05-06T08:00:00Z")

	root, g := cli.NewRoot()
	root.AddCommand(NewVersion(g))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"version", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var got version.Info
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v\nbody=%s", err, buf.String())
	}
	if got.Version != "1.2.3" || got.Commit != "abc123" || got.Date != "2026-05-06T08:00:00Z" {
		t.Fatalf("unexpected version info: %+v", got)
	}
}
