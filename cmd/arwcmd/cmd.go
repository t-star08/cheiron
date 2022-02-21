package arwcmd

import (
	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/cmd/arwcmd/multicmd"
	"github.com/t-star08/cheiron/cmd/arwcmd/routinecmd"
	"github.com/t-star08/cheiron/cmd/arwcmd/singlecmd"
)

var CMD = &cobra.Command{
	Use: "arrow",
}

func init() {
	CMD.AddCommand (
		routinecmd.CMD,
		singlecmd.CMD,
		multicmd.CMD,
	)
}
