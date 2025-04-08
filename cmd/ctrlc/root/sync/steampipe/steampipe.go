package steampipe

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/ctrlplanedev/cli/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewSyncSteampipeCmd() *cobra.Command {
	var providerName string
	var spConnection string
	var spTable string

	apiURL := viper.GetString("url")
	apiKey := viper.GetString("api-key")
	workspaceId := viper.GetString("workspace")

	cmd := &cobra.Command{
		Use:   "steampipe",
		Short: "Subcommands for integrating steampipe with Ctrlplane",
		Example: heredoc.Doc(`
			$ ctrlc sync steampipe -r resource-provider -c steampipe-connection -t steampipe-table
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var provider *api.ResourceProvider
			var providerResp *http.Response

			ctx := context.Background()

			spClient, err := NewSteampipeClient(spConnection)
			if err != nil {
				return fmt.Errorf("Failed to create Steampipe Connection: %w", err)
			}

			cpClient, err := api.NewAPIKeyClientWithResponses(apiURL, apiKey)

			if err != nil {
				return fmt.Errorf("failed to create API spClient: %w", err)
			}

			defer spClient.Close()

			resources, err := spClient.Fetch(spTable)
			if err != nil {
				return err
			}

			log.Infof("Resource count  %d", len(resources))
			if len(resources) > 0 {
				provider, err = api.NewResourceProvider(cpClient, workspaceId, providerName)
				if err != nil {
					return fmt.Errorf("failed to create resource provider: %w", err)
				}
				providerResp, err = provider.UpsertResource(ctx, resources)
				if providerResp != nil {
					log.Info("upsert resources response ", "status", providerResp.Status)
				}
				if err != nil {
					return fmt.Errorf("failed to upsert resources: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&providerName, "resource-provider-name", "r", os.Getenv("RESOURCE_PROVIDER"), "The resource group name")
	cmd.Flags().StringVarP(&spConnection, "steampipe-connection", "c", os.Getenv("STEAMPIPE_CONNECTION"), "The steampipe postgresql connection string to use")
	cmd.Flags().StringVarP(&spTable, "steampipe-table", "t", os.Getenv("STEAMPIPE_TABLE"), "The steampipe postgresql table to select from")

	cmd.MarkFlagRequired("resource-provider-name")
	cmd.MarkFlagRequired("steampipe-connection")
	cmd.MarkFlagRequired("steampipe-table")

	return cmd
}
