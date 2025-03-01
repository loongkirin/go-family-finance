package request

type QueryWhere struct {
	QueryFilters []*QueryFilter `json:"query_filters"`
	Connector    Connector      `json:"connector"`
}

func NewQueryWhere(filters []*QueryFilter, connector Connector) *QueryWhere {
	return &QueryWhere{
		QueryFilters: filters,
		Connector:    connector,
	}
}
