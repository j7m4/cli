package adapter

import "github.com/ctrlplanedev/cli/internal/api"

var awsAdapters = []SteampipeAdapter{
	{
		Table: "aws_eks_cluster",
		Translate: func(data *map[string]interface{}) (api.AgentResource, bool) {
			var ok bool
			var resource api.AgentResource
			var tags map[string]string

			resource = api.AgentResource{
				Config:   make(map[string]interface{}),
				Metadata: make(map[string]string),
			}

			if resource.Identifier, ok = getValue[string](data, []string{"arn"}); !ok {
				return resource, false
			}

			if resource.Name, ok = getValue[string](data, []string{"name"}); !ok {
				return resource, false
			}

			if tags, ok = getValue[map[string]string](data, []string{"tags"}); ok {
				for key, value := range tags {
					resource.Metadata[key] = value
				}
			} else {
				return resource, false
			}

			if resource.Version, ok = getValue[string](data, []string{"version"}); !ok {
				return resource, false
			}

			if resource.Kind, ok = getValue[string](data, []string{"kind"}); !ok {
				return resource, false
			}

			return resource, true
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
