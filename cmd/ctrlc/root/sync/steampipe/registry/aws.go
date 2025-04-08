package registry

var awsComponents = []SteampipeAccessInfo{
	{
		ConnectionType: "aws",
		/*
			IsType: func(resourceType string) bool {
				return resourceType == "eks_cluster"
			},
			Convert func(res) ResourceObj {
				return ResourceObj{
					Name : fooo,
					Version: "kubernetes/v1",
					Kind: "ClusterApi",
					Config: map[string]interface{}{
						"auth": map[string]interface{} {
							method: "aws/eks",
							region: "us-east-1",
							cluster: "default",
							accountId: "123456789012",
						},
						status: "active",
						name: "default",
					},
					Identifier: "my-arn",
					Metadata: map[string]interface{}{
						"bunch": "of",
						"misc": "stuff",
						"not": "excludeed",
					},

			},
		*/
		ResourceType: "eks_cluster",
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
