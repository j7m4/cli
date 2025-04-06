package steampipe

import (
	"fmt"
	"sort"

	"github.com/ctrlplanedev/cli/cmd/ctrlc/root/sync/steampipe/registry"
)

type Connection struct {
	Name         string
	Type         string
	ResourceType string
}

// ForeignTable represents a row from information_schema.foreign_tables
type ForeignTable struct {
	ForeignTableCatalog string
	ForeignTableSchema  string
	ForeignTableName    string
	ForeignServerName   string
}

func (c *SteampipeClient) ListResourceGroups() ([]Connection, error) {
	resourceGroups := make([]Connection, 0)

	foreignTables, err := c.GetSteampipeForeignTables()
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign tables: %w", err)
	}

	for _, foreignTable := range foreignTables {
		component, ok := registry.GetAccessInfo(foreignTable.ForeignTableName)
		if !ok {
			continue
		}
		connectionName := foreignTable.ForeignTableSchema
		if component.ConnectionType == connectionName {
			connectionName = "*"
		}
		resourceGroups = append(resourceGroups, Connection{
			Name:         connectionName,
			Type:         component.ConnectionType,
			ResourceType: component.ResourceType,
		})
	}

	fmt.Printf("%-30s %-20s %-30s\n", "connection-name", "connection-type", "resource-type")
	fmt.Printf("%-30s %-20s %-30s\n", "------------------------------", "--------------------", "--------------------------------")
	sort.Slice(resourceGroups, func(i, j int) bool {
		if resourceGroups[i].Name == resourceGroups[j].Name {
			if resourceGroups[i].Type == resourceGroups[j].Type {
				return resourceGroups[i].ResourceType < resourceGroups[j].ResourceType
			}
			return resourceGroups[i].Type < resourceGroups[j].Type
		}
		return resourceGroups[i].Name < resourceGroups[j].Name
	})
	for _, rg := range resourceGroups {
		fmt.Printf("%-30s %-20s %-30s\n", rg.Name, rg.Type, rg.ResourceType)
	}

	return resourceGroups, nil
}

// GetSteampipeForeignTables returns all foreign tables where foreign_table_catalog is 'steampipe'
func (c *SteampipeClient) GetSteampipeForeignTables() ([]ForeignTable, error) {
	query := `
		SELECT 
			foreign_table_catalog,
			foreign_table_schema,
			foreign_table_name,
			foreign_server_name
		FROM information_schema.foreign_tables
		WHERE foreign_table_catalog = 'steampipe'
		ORDER BY foreign_table_schema, foreign_table_name
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign tables: %w", err)
	}
	defer rows.Close()

	var tables []ForeignTable
	for rows.Next() {
		var table ForeignTable
		err := rows.Scan(
			&table.ForeignTableCatalog,
			&table.ForeignTableSchema,
			&table.ForeignTableName,
			&table.ForeignServerName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan foreign table row: %w", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foreign table rows: %w", err)
	}

	return tables, nil
}
