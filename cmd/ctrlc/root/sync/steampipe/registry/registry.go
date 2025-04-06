package registry

func buildRegistry() []SteampipeAccessInfo {
	components := []SteampipeAccessInfo{}
	components = append(components, kubernetesComponents...)
	components = append(components, awsComponents...)
	components = append(components, gcpComponents...)
	components = append(components, azureComponents...)
	return components
}

func buildTableNameRegistry() map[string]SteampipeAccessInfo {
	tableNameRegistry := make(map[string]SteampipeAccessInfo)
	for _, component := range buildRegistry() {
		tableNameRegistry[component.TableName()] = component
	}
	return tableNameRegistry
}

var tableNameRegistry = buildTableNameRegistry()

func GetAccessInfo(tableName string) (SteampipeAccessInfo, bool) {
	accessInfo, ok := tableNameRegistry[tableName]
	return accessInfo, ok
}
