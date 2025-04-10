package model

type SteampipeAdapterStruct struct {
	Table   string
	Convert func(jsonRow string) (string, bool)
}

func (a *SteampipeAdapterStruct) EntityName() string {
	return a.Table
}

func (a *SteampipeAdapterStruct) ToResourceJson(sqlRow SqlRow) (string, bool) {
	return a.Convert(sqlRow.Json)
}

type SteampipeAdapter interface {
	EntityName() string
	ToResourceJson(sqlRow SqlRow) (string, bool)
}
