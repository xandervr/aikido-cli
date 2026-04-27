package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewResearch(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "research", Short: "Vulnerability research lookups"}
	cmd.AddCommand(
		simpleGet(g, "cve <cve-id>", "Get CVE details", "/research/cves"),
		simpleGet(g, "changelog <package>", "Package changelog summary", "/research/changelogs"),
		simpleList(g, "malware-packages", "List recently flagged malware packages", "/research/malware-packages"),
	)
	return cmd
}

func NewCVE(g *cli.Globals) *cobra.Command {
	return simpleGet(g, "cve <cve-id>", "Get CVE details (shortcut)", "/research/cves")
}

func NewChangelog(g *cli.Globals) *cobra.Command {
	return simpleGet(g, "changelog <package>", "Package changelog (shortcut)", "/research/changelogs")
}

func NewMalwarePackages(g *cli.Globals) *cobra.Command {
	return simpleList(g, "malware-packages", "Malware packages (shortcut)", "/research/malware-packages")
}
