package model

import (
	"github.com/ctrlplanedev/cli/internal/api"
)

type SteampipeAdapterStruct struct {
	Table     string
	Translate func(data map[string]interface{}) (api.AgentResource, bool)
}

func (a *SteampipeAdapterStruct) EntityName() string {
	return a.Table
}

func (a *SteampipeAdapterStruct) ToApiResource(row SqlRow) (api.AgentResource, bool) {
	return a.Translate(row.Data)
}

type SteampipeAdapter interface {
	EntityName() string
	ToApiResource(row SqlRow) (api.AgentResource, bool)
}
