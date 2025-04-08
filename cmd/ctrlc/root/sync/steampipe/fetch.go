package steampipe

import (
	"fmt"
	"github.com/ctrlplanedev/cli/internal/api"
)

func (c *SteampipeClient) Fetch(tableName string) ([]api.AgentResource, error) {
	resources := make([]api.AgentResource, 0)

	results, err := c.SelectAll(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from %s table: %w", tableName, err)
	}

	for _, result := range results {

		fmt.Println(result)

		resources = append(resources, api.AgentResource{})
	}

	return resources, nil
}

// SelectAll returns all foreign tables where foreign_table_catalog is 'steampipe'
func (c *SteampipeClient) SelectAll(tableName string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s table: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var row = make(map[string]interface{})
		err := rows.Scan(row)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s table row: %w", tableName, err)
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating %s table rows: %w", tableName, err)
	}

	return results, nil
}
