package steampipe

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ctrlplanedev/cli/cmd/ctrlc/root/sync/steampipe/adapter"
	"github.com/ctrlplanedev/cli/internal/api"
)

func (c *SteampipeClient) Fetch(table string) ([]api.AgentResource, error) {
	var resource api.AgentResource
	var ok bool

	ad := adapter.SelectAdapter(table)
	if ad == nil {
		return nil, fmt.Errorf("no adapter found for table %s", table)
	}

	resources := make([]api.AgentResource, 0)

	results, err := c.SelectAll(table)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from %s table: %w", table, err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	for _, result := range results {

		fmt.Println(result)

		if resource, ok = ad.Translate(result); ok {
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// SelectAll returns all foreign tables where foreign_table_catalog is 'steampipe'
func (c *SteampipeClient) SelectAll(tableName string) ([]*map[string]interface{}, error) {
	var columns []string
	var err error

	results := make([]*map[string]interface{}, 0)

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s table: %w", tableName, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorf("failed to close rows: %v", err)
		}
	}(rows)

	columns, err = rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns from %s table: %w", tableName, err)
	}

	var row map[string]interface{}
	for rows.Next() {
		row, err = toJsonObj(rows, columns)
		results = append(results, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating %s table rows: %w", tableName, err)
	}

	return results, nil
}

// toJsonObj converts a row from the database to a JSON object
// It assumes that a next row is available!
func toJsonObj(nextRow *sql.Rows, columns []string) (map[string]interface{}, error) {
	// Create values and associate ptrs to them for row scanning
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	jsonObj := make(map[string]interface{})
	jsonArr := make([]interface{}, 0)

	err := nextRow.Scan(valuePtrs...)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	result := make(map[string]interface{})
	for i, col := range columns {
		if values[i] == nil {
			continue
		}
		switch v := values[i].(type) {
		case []byte:
			if err = json.Unmarshal(v, &jsonObj); err == nil {
				values[i] = jsonObj
				break
			}
			if err = json.Unmarshal(v, &jsonArr); err == nil {
				values[i] = jsonArr
				break
			}
			values[i] = string(v)
		default:
			values[i] = v
		}
		result[col] = values[i]
	}
	return result, nil
}
