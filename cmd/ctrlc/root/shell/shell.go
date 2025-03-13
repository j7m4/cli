package shell

import (
	"fmt"
	"os"

	"github.com/ctrlplanedev/cli/internal/repl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Interact with ctrlplane via shell interface.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return repl.StartLoop(cmd)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			apiURL := viper.GetString("url")
			if apiURL == "" {
				fmt.Fprintln(cmd.ErrOrStderr(), "API URL is required. Set via --url flag or in config")
				os.Exit(1)
			}
			apiKey := viper.GetString("api-key")
			if apiKey == "" {
				fmt.Fprintln(cmd.ErrOrStderr(), "API key is required. Set via --api-key flag or in config")
				os.Exit(1)
			}
			return nil
		},
	}

	return cmd
}
