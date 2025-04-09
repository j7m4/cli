package adapter

import "github.com/ctrlplanedev/cli/internal/api"

var awsAdapters = []SteampipeAdapter{
	{
		Table: "aws_eks_cluster",
		Translate: func(data *map[string]interface{}) (api.AgentResource, bool) {
			var name string
			var clusterArn string
			var tags map[string]string
			var accountId string
			var region string
			var certificateAuthorityData string
			var status string
			var version string
			var endpoint string
			var roleArn string
			var platformVersion string

			var zero api.AgentResource = api.AgentResource{}
			var ok bool

			if clusterArn, ok = getPathValue[string](data, "arn"); !ok {
				return zero, false
			}
			if name, ok = getPathValue[string](data, "name"); !ok {
				return zero, false
			}
			if tags, ok = getPathValue[map[string]string](data, "tags"); !ok {
				return zero, false
			}
			if accountId, ok = getPathValue[string](data, "account_id"); !ok {
				return zero, false
			}
			if region, ok = getPathValue[string](data, "region"); !ok {
				return zero, false
			}
			if certificateAuthorityData, ok = getPathValue[string](data, "certificate_authority.Data"); !ok {
				return zero, false
			}
			if status, ok = getPathValue[string](data, "status"); !ok {
				return zero, false
			}
			if version, ok = getPathValue[string](data, "version"); !ok {
				return zero, false
			}
			if endpoint, ok = getPathValue[string](data, "endpoint"); !ok {
				return zero, false
			}
			if roleArn, ok = getPathValue[string](data, "role_arn"); !ok {
				return zero, false
			}
			if platformVersion, ok = getPathValue[string](data, "platform_version"); !ok {
				return zero, false
			}

			metadata := map[string]string{
				"aws/account-id":       accountId,
				"aws/arn":              clusterArn,
				"aws/eks-role-arn":     roleArn,
				"aws/region":           region,
				"aws/platform-version": platformVersion,
			}
			for k, v := range tags {
				metadata[k] = v
			}

			return api.AgentResource{
				Identifier: clusterArn,
				Name:       name,
				Version:    version,
				Kind:       "ClusterAPI",
				Config: map[string]interface{}{
					"auth": map[string]string{
						"method":      "aws/eks",
						"region":      region,
						"accountId":   accountId,
						"clusterName": name,
					},
					"name": name,
					"server": map[string]string{
						"endpoint":                   endpoint,
						"certificationAuthorityData": certificateAuthorityData,
					},
					"status": status,
				},
				Metadata: metadata,
			}, true
		},
	},
	//{
	//	Plugin:       "aws",
	//	ResourceType: "vpc",
	//},
	//{
	//	Plugin:       "aws",
	//	ResourceType: "rds_db_cluster",
	//},
	//{
	//	Plugin:       "aws",
	//	ResourceType: "rds_db_instance",
	//},
	//{
	//	Plugin:       "aws",
	//	ResourceType: "elasticache_cluster",
	//},
}
