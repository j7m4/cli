package registry

var kubernetesComponents = []SteampipeAccessInfo{
	{
		ConnectionType: "kubernetes",
		ResourceType:   "namespace",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "pod",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "node",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "service",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "deployment",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "job",
	},
	{
		ConnectionType: "kubernetes",
		ResourceType:   "cronjob",
	},
}
