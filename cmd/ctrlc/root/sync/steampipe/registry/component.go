package registry

import "fmt"

type SteampipeAccessInfo struct {
	ConnectionType string
	ResourceType   string
	TablePrefix    string
}

func (c *SteampipeAccessInfo) GetTablePrefix() string {
	if c.TablePrefix == "" {
		return c.ConnectionType
	}
	return c.TablePrefix
}

func (c *SteampipeAccessInfo) TableName() string {
	return fmt.Sprintf("%s_%s", c.GetTablePrefix(), c.ResourceType)
}
