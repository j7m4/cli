package adapter

func buildRegistry() []SteampipeAdapter {
	components := []SteampipeAdapter{}
	components = append(components, kubernetesComponents...)
	components = append(components, awsAdapters...)
	components = append(components, gcpComponents...)
	components = append(components, azureComponents...)
	return components
}

func buildTableNameRegistry() map[string]SteampipeAdapter {
	tableNameRegistry := make(map[string]SteampipeAdapter)
	for _, component := range buildRegistry() {
		tableNameRegistry[component.TableName()] = component
	}
	return tableNameRegistry
}

var tableNameRegistry = buildTableNameRegistry()

func GetAccessInfo(tableName string) (SteampipeAdapter, bool) {
	accessInfo, ok := tableNameRegistry[tableName]
	return accessInfo, ok
}
