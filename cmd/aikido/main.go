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
	)
	if err := root.Execute(); err != nil {
		cli.Exit(err)
	}
}
