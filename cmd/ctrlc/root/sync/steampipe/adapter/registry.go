package adapter

import (
	"github.com/charmbracelet/log"
	"strings"
)

func buildRegistry() map[string]*SteampipeAdapter {
	registry := make(map[string]*SteampipeAdapter)
	for _, pluginAdapters := range [][]SteampipeAdapter{kubernetesComponents, awsAdapters, gcpComponents, azureComponents} {
		for _, adapter := range pluginAdapters {
			registry[adapter.Table] = &adapter
		}
	}
	return registry
}

var registry = buildRegistry()

func SelectAdapter(table string) *SteampipeAdapter {
	var adapter *SteampipeAdapter
	var ok bool
	tableName := stripSchema(table)

	if adapter, ok = registry[tableName]; !ok {
		log.Warnf("could not find adapter for table %s", tableName)
		return nil
	}
	return adapter
}

func stripSchema(table string) string {
	if strings.Contains(table, ".") {
		return strings.Split(table, ".")[1]
	}
	return table
}
