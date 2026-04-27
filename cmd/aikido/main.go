package main

import (
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/cli/commands"
)

func main() {
	root, g := cli.NewRoot()
	root.AddCommand(commands.NewAuth(g))
	if err := root.Execute(); err != nil {
		cli.Exit(err)
	}
}
