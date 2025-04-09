package adapter

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/ctrlplanedev/cli/internal/api"
)

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
			var err error
			var entityName = "aws_eks_cluster"

			if clusterArn, err = getPathValue[string](data, "arn"); err != nil {
				log.Infof("could not find %s in %s: %s", "arn", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if name, err = getPathValue[string](data, "name"); err != nil {
				log.Infof("could not find %s in %s: %s", "name", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if tags, err = getPathValue[map[string]string](data, "tags"); err != nil {
				log.Infof("could not find %s in %s: %s", "tags", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if accountId, err = getPathValue[string](data, "account_id"); err != nil {
				log.Infof("could not find %s in %s: %s", "account_id", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if region, err = getPathValue[string](data, "region"); err != nil {
				log.Infof("could not find %s in %s: %s", "region", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if certificateAuthorityData, err = getPathValue[string](data, "certificate_authority.Data"); err != nil {
				log.Infof("could not find %s in %s: %s", "certificate_authority.Data", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if status, err = getPathValue[string](data, "status"); err != nil {
				log.Infof("could not find %s in %s: %s", "status", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if version, err = getPathValue[string](data, "version"); err != nil {
				log.Infof("could not find %s in %s: %s", "version", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if endpoint, err = getPathValue[string](data, "endpoint"); err != nil {
				log.Infof("could not find %s in %s: %s", "endpoint", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if roleArn, err = getPathValue[string](data, "role_arn"); err != nil {
				log.Infof("could not find %s in %s: %s", "role_arn", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
				return zero, false
			}
			if platformVersion, err = getPathValue[string](data, "platform_version"); err != nil {
				log.Infof("could not find %s in %s: %s", "platform_version", entityName, err.Error())
				dataStr, _ := json.Marshal(data)
				log.Debugf("data: %s", dataStr)
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
