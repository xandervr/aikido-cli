package main

import (
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/cli/commands"
)

func main() {
	root, g := cli.NewRoot()
	root.AddCommand(
		commands.NewAuth(g),
		commands.NewWorkspace(g),
		commands.NewRepos(g),
		commands.NewIssues(g),
		commands.NewTeams(g),
		commands.NewUsers(g),
		commands.NewContainers(g),
		commands.NewClouds(g),
		commands.NewApps(g),
		commands.NewVMs(g),
		commands.NewLicenses(g),
		commands.NewWebhooks(g),
		commands.NewActivity(g),
		commands.NewPRChecks(g),
		commands.NewCompliance(g),
		commands.NewCustomRules(g),
		commands.NewPentest(g),
		commands.NewTasks(g),
		commands.NewResearch(g),
		commands.NewCVE(g),
		commands.NewChangelog(g),
		commands.NewMalwarePackages(g),
		commands.NewReport(g),
		commands.NewVersion(g),
	)
	if err := root.Execute(); err != nil {
		cli.Exit(err)
	}
}
