package adapter

import (
	"fmt"
	"github.com/ctrlplanedev/cli/internal/api"
)

type SteampipeAdapter struct {
	ConnectionType string
	ResourceType   string
	TablePrefix    string
	Translate      func(data map[string]interface{}) api.AgentResource
	IsCompatible   func(data map[string]interface{}) bool
}

func (c *SteampipeAdapter) GetTablePrefix() string {
	if c.TablePrefix == "" {
		return c.ConnectionType
	}
	return c.TablePrefix
}

func (c *SteampipeAdapter) TableName() string {
	return fmt.Sprintf("%s_%s", c.GetTablePrefix(), c.ResourceType)
}
