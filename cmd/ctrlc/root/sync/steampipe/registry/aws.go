package registry

var awsComponents = []SteampipeAccessInfo{
	{
		ConnectionType: "aws",
		ResourceType:   "eks_cluster",
	},
	{
		ConnectionType: "aws",
		ResourceType:   "vpc",
	},
	{
		ConnectionType: "aws",
		ResourceType:   "rds_db_cluster",
	},
	{
		ConnectionType: "aws",
		ResourceType:   "rds_db_instance",
	},
	{
		ConnectionType: "aws",
		ResourceType:   "elasticache_cluster",
	},
}
