package shazam

import (
	"github.com/electric-saw/pg-shazam/pkg/cmd/completion"
	"github.com/electric-saw/pg-shazam/pkg/cmd/start"
	"github.com/electric-saw/pg-shazam/pkg/cmd/version"
	"github.com/spf13/cobra"
)

func NewShazamCommand(appName string) *cobra.Command {
	root := &cobra.Command{
		Use:   appName,
		Short: "Pg-shazam",
	}

	root.AddCommand(version.NewCmdVersion(appName))
	root.AddCommand(completion.NewCmdCompletion())
	root.AddCommand(start.NewCmdStart())

	return root

}
