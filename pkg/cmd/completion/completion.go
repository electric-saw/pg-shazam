package completion

import (
	"io"
	"os"

	"github.com/electric-saw/pg-shazam/pkg/util"
	"github.com/spf13/cobra"
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"powershell": runCompletionPowerShell,
	}
)

func NewCmdCompletion() *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	return &cobra.Command{
		Use:    "completion SHELL",
		Short:  "Output shell completion code for the specified shell (bash, zsh or powershell)",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := runCompletion(os.Stdout, cmd, args)
			util.CheckErr(err)
		},
		ValidArgs: shells,
	}
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return util.UsageErrorf(cmd, "Shell not specified.")
	}
	if len(args) > 1 {
		return util.UsageErrorf(cmd, "Too many arguments. Expected only the shell type.")
	}
	run, found := completionShells[args[0]]
	if !found {
		return util.UsageErrorf(cmd, "Unsupported shell type %q.", args[0])
	}

	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, root *cobra.Command) error {
	return root.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, root *cobra.Command) error {
	return root.GenZshCompletion(out)
}

func runCompletionPowerShell(out io.Writer, root *cobra.Command) error {
	return root.GenPowerShellCompletion(out)
}
