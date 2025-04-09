package aws

import (
	"github.com/ctrlplanedev/cli/cmd/ctrlc/root/sync/steampipe/adapter/model"
	"github.com/ctrlplanedev/cli/internal/api"
)

const ec2Table = "aws_ec2_instance"

var EC2 model.SteampipeAdapter = &model.SteampipeAdapterStruct{
	Table: ec2Table,
	Translate: func(data map[string]interface{}) (api.AgentResource, bool) {
		var entityName = ec2Table
		var sqlRow model.SqlRow = model.SqlRow{
			EntityName: entityName,
			Data:       data,
		}

		var name string
		var clusterArn string
		var tags model.Tags
		var accountId string
		var region string
		var status string
		var endpoint string

		var zero api.AgentResource = api.AgentResource{}

		if !model.GetRequiredValue[string](sqlRow, "arn", &clusterArn) ||
			!model.GetRequiredValue[string](sqlRow, "instance_id", &name) ||
			!model.GetOptionalValue[model.Tags](sqlRow, "tags", &tags) ||
			!model.GetRequiredValue[string](sqlRow, "account_id", &accountId) ||
			!model.GetRequiredValue[string](sqlRow, "region", &region) ||
			!model.GetRequiredValue[string](sqlRow, "instance_state", &status) ||
			!model.GetRequiredValue[string](sqlRow, "private_dns_name", &endpoint) {
			return zero, false
		}

		metadata := model.BuildMetadata(map[string]string{
			"aws/account-id": accountId,
			"aws/arn":        clusterArn,
			"aws/region":     region,
		}).AppendTags(tags)

		return api.AgentResource{
			Identifier: clusterArn,
			Name:       name,
			Version:    "compute/v1",
			Kind:       "Compute",
			Config: map[string]interface{}{
				"auth": map[string]string{
					"method":     "aws/ec2",
					"region":     region,
					"accountId":  accountId,
					"instanceId": name,
				},
				"name": name,
				"server": map[string]string{
					"endpoint": endpoint,
				},
				"status": status,
			},
			Metadata: metadata,
		}, true
	},
}
