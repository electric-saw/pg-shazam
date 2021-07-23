package start

import (
	"github.com/electric-saw/pg-shazam/internal/pkg/frontend"
	"github.com/electric-saw/pg-shazam/pkg/util"
	"github.com/spf13/cobra"
)

func NewCmdStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start <config file>",
		Short: "Start frontend server to pg",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doRun(args[0])
		},
	}
}

func doRun(file string) {
	f, err := frontend.NewFrontend(file)
	util.CheckErr(err)

	defer f.Close()

	util.CheckErr(f.Run())
}
