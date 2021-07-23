package version

import (
	"fmt"

	"github.com/electric-saw/pg-shazam/internal/version"
	"github.com/spf13/cobra"
)

func NewCmdVersion(appName string) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   fmt.Sprintf("Print the %s version", appName),
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}

func printVersion() {
	v := version.Get()
	fmt.Println("Version        ", v.Version)
	fmt.Println("Git commit     ", v.GitCommit)
	fmt.Println("Git tree state ", v.GitTreeState)
	fmt.Println("Go version     ", v.GoVersion)
}
