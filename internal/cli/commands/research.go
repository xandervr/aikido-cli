package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewResearch(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "research", Short: "Vulnerability research lookups"}
	cmd.AddCommand(
		simpleGet(g, "cve <cve-id>", "Get CVE details", "/cve"),
		changelogCommand(g, "changelog <package>", "Package changelog summary"),
		malwarePackagesCommand(g, "malware-packages", "List recently flagged malware packages"),
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
	return malwarePackagesCommand(g, "malware-packages", "Malware packages (shortcut)")
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

func malwarePackagesCommand(g *cli.Globals, use, short string) *cobra.Command {
	var page, perPage int
	var search, ecosystem string
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if page > 0 {
				q["page"] = fmt.Sprintf("%d", page)
			}
			if perPage > 0 {
				q["per_page"] = fmt.Sprintf("%d", perPage)
			}
			if search != "" {
				q["search"] = search
			}
			if ecosystem != "" {
				q["filter_ecosystem"] = ecosystem
			}
			var raw any
			if err := c.Get(cmd.Context(), "/research/malware/packages", q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page (0-indexed)")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "page size (10-20, default 20)")
	cmd.Flags().StringVar(&search, "search", "", "search packages")
	cmd.Flags().StringVar(&ecosystem, "ecosystem", "", "filter by ecosystem: npm|cargo|pypi|packagist|golang|maven|nuget|ruby|open_vsx|github_act|vscode|chrome|wordpress|skills_sh")
	return cmd
}
