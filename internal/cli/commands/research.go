package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewResearch(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "research", Short: "Vulnerability research lookups"}
	cmd.AddCommand(
		simpleGet(g, "cve <cve-id>", "Get CVE details", "/cve"),
		changelogCommand(g, "changelog <package>", "Package changelog summary"),
		simpleList(g, "malware-packages", "List recently flagged malware packages", "/research/malware/packages"),
	)
	return cmd
}

func NewCVE(g *cli.Globals) *cobra.Command {
	return simpleGet(g, "cve <cve-id>", "Get CVE details (shortcut)", "/cve")
}

func NewChangelog(g *cli.Globals) *cobra.Command {
	return changelogCommand(g, "changelog <package>", "Package changelog (shortcut)")
}

func NewMalwarePackages(g *cli.Globals) *cobra.Command {
	return simpleList(g, "malware-packages", "Malware packages (shortcut)", "/research/malware/packages")
}

func changelogCommand(g *cli.Globals, use, short string) *cobra.Command {
	var fromVersion, toVersion, language string
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromVersion == "" || toVersion == "" || language == "" {
				return errors.New("--from, --to, and --language are required")
			}
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{
				"package_name": args[0],
				"from_version": fromVersion,
				"to_version":   toVersion,
				"language":     language,
			}
			var raw any
			if err := c.Get(cmd.Context(), "/changelog-summary", q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	cmd.Flags().StringVar(&fromVersion, "from", "", "current package version (required)")
	cmd.Flags().StringVar(&toVersion, "to", "", "target package version (required)")
	cmd.Flags().StringVar(&language, "language", "", "package language: JS|PY|GO|.NET|Java|Scala|Kotlin (required)")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("language")
	return cmd
}
