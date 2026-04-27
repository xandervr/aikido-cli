package main

import (
	"github.com/xandervr/aikido-cli/internal/cli"
)

func main() {
	root, _ := cli.NewRoot()
	if err := root.Execute(); err != nil {
		cli.Exit(err)
	}
}
