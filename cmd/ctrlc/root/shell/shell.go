package shell

import (
	"github.com/ctrlplanedev/cli/internal/repl"
	"github.com/spf13/cobra"
)

func NewShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Interact with ctrlplane via shell interface.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return repl.StartLoop()
		},
	}

	return cmd
}
