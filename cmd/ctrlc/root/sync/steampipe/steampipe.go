package steampipe

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/ctrlplanedev/cli/internal/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewSyncSteampipeCmd() *cobra.Command {
	apiURL := viper.GetString("url")
	apiKey := viper.GetString("api-key")
	workspaceId := viper.GetString("workspace")

	ctrlplaneClient, err := api.NewAPIKeyClientWithResponses(apiURL, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	cmd := &cobra.Command{
		Use:   "steampipe <subcommand>",
		Short: "Subcommands for integrating steampipe with Ctrlplane",
		Example: heredoc.Doc(`
			$ ctrlc sync steampipe list                         # Show which resource providers are available
			$ ctrlc sync steampipe send <resource-provider-id>  # Send to Ctrlplane the resource info for all resourceGroups
		`),
	}

	cmd.AddCommand(NewSyncSteampipeListCmd())
	cmd.AddCommand(NewSyncSteampipeSendCmd())

	return cmd
}

func NewSyncSteampipeListCmd() *cobra.Command {
	var spConnection string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available resource-providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Listing all available resourceGroups")

			client, err := NewSteampipeClient(spConnection)
			if err != nil {
				log.Error("Failed to create Steampipe client", "error", err)
				return err
			}
			defer client.Close()

			connections, err := client.ListConnections()
			if err != nil {
				return err
			}

			// Create a new table writer
			table := tablewriter.NewWriter(os.Stdout)

			// Set headers
			table.SetHeader([]string{"Resource ID", "Resource Type", "Connection Name", "Steampipe Table"})

			// Add rows
			for _, conn := range connections {
				table.Append([]string{
					conn.CtrlPlaneResource.Id,
					conn.CtrlPlaneResource.Type,
					conn.Name,
					conn.SteampipeResource.TableName,
				})
			}

			// Set table properties
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("|")
			table.SetRowSeparator("")
			table.SetBorder(true)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)

			// Render the table
			table.Render()

			return nil
		},
	}

	return cmd
}

func NewSyncSteampipeSendCmd() *cobra.Command {
	var resourceGroup string
	var spConnection string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send resource info to Ctrlplane",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if resourceGroup == "" {
				return fmt.Errorf("resource-group must be provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Sending resource groups to Ctrlplane")

			client, err := NewSteampipeClient(spConnection)
			if err != nil {
				log.Error("Failed to create Steampipe client", "error", err)
				return err
			}
			defer client.Close()

			resourceGroups, err := client.SendResourcesFromGroup(resourceGroup)
			if err != nil {
				return err
			}

			for _, group := range resourceGroups {
				fmt.Println(group)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&resourceGroup, "resource-group", "r", os.Getenv("STEAMPIPE_RESOURCE_GROUP"), "The resource group name")
	cmd.Flags().StringVarP(&spConnection, "steampipe-connection", "c", "", "The steampipe postgresql connection string to use")

	cmd.MarkFlagRequired("resource-group")
	cmd.MarkFlagRequired("steampipe-connection")

	return cmd
}
