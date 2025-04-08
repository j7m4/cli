package steampipe

import (
	"fmt"
	"sort"

	"github.com/ctrlplanedev/cli/cmd/ctrlc/root/sync/steampipe/registry"
)

type Connection struct {
	Name              string
	Type              string
	SteampipeResource SteampipeResource
	CtrlPlaneResource CtrlPlaneResource
}

type CtrlPlaneResource struct {
	Id   string
	Type string
}

type SteampipeResource struct {
	TableName      string
	ConnectionName string
}

// ForeignTable represents a row from information_schema.foreign_tables
type ForeignTable struct {
	ForeignTableCatalog string
	ForeignTableSchema  string
	ForeignTableName    string
	ForeignServerName   string
}

// SortedConnections implements sort.Interface for []Connection based on Name, Type, and CtrlPlaneResource.Type
type SortedConnections []Connection

func (c SortedConnections) Len() int      { return len(c) }
func (c SortedConnections) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c SortedConnections) Less(i, j int) bool {
	if c[i].Name == c[j].Name {
		if c[i].Type == c[j].Type {
			return c[i].CtrlPlaneResource.Type < c[j].CtrlPlaneResource.Type
		}
		return c[i].Type < c[j].Type
	}
	return c[i].Name < c[j].Name
}

func (c *SteampipeClient) ListConnections() ([]Connection, error) {
	connections := make([]Connection, 0)

	foreignTables, err := c.GetSteampipeForeignTables()
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign tables: %w", err)
	}

	for _, foreignTable := range foreignTables {

		spResource := SteampipeResource{
			TableName:      foreignTable.ForeignTableName,
			ConnectionName: foreignTable.ForeignTableSchema,
		}

		component, ok := registry.GetAccessInfo(foreignTable.ForeignTableName)
		if !ok {
			continue
		}

		cpResource := CtrlPlaneResource{
			Id:   "",
			Type: component.ResourceType,
		}

		// Connection name defaults to the schema name, but for the case
		// where the schema name is the same as the .... TODO: fix this descr
		name := foreignTable.ForeignTableSchema
		if component.ConnectionType == name {
			name = "*"
		}

		conn := Connection{
			Name:              name,
			Type:              component.ConnectionType,
			SteampipeResource: spResource,
			CtrlPlaneResource: cpResource,
		}

		connections = append(connections, conn)
	}

	// Replace the sort.Slice with sort.Sort using our custom sort.Interface
	sort.Sort(SortedConnections(connections))

	return connections, nil
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
