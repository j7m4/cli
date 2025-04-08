package adapter

var gcpComponents = []SteampipeAdapter{
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
