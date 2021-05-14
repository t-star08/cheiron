package project

import (
	"github.com/t-star08/cheiron/cmd/project/arwcmd"
	"github.com/t-star08/cheiron/cmd/project/initcmd"
	"github.com/t-star08/cheiron/cmd/project/qvrcmd"

	"github.com/spf13/cobra"
)

var CMD = &cobra.Command {
	Use: "project",
}

func init() {
	CMD.AddCommand (
		arwcmd.CMD,
		initcmd.CMD,
		qvrcmd.CMD,
	)
}
