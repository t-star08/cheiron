package main

import (
	"github.com/t-star08/cheiron/cmd/insert/istcmd"
	"github.com/t-star08/cheiron/cmd/project"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command {
	Use: "cheiron",
	Version: "v0.0.0",
}

func init() {
	cmd.AddCommand (
		project.CMD,
		istcmd.CMD,
	)
}

func main() {
	cmd.Execute()
}
