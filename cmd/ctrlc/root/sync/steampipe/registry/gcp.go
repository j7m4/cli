package registry

var gcpComponents = []SteampipeAccessInfo{
	{
		ConnectionType: "gcp",
		ResourceType:   "kubernetes_cluster",
	},
	{
		ConnectionType: "gcp",
		ResourceType:   "sql_database",
	},
	{
		ConnectionType: "gcp",
		ResourceType:   "sql_database_instance",
	},
}
