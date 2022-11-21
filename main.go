package main

import (
	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/cmd/arwcmd"
	"github.com/t-star08/cheiron/cmd/initcmd"
	"github.com/t-star08/cheiron/cmd/statuscmd"
)

var cmd = &cobra.Command {
	Use: "cheiron",
	Version: "v2.0.1",
}

func init() {
	cmd.AddCommand (
		arwcmd.CMD,
		initcmd.CMD,
		statuscmd.CMD,
	)
}

func main() {
	cmd.Execute()
}
