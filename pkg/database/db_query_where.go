package database

type DbQueryWhere struct {
	QueryFilters []DbQueryFilter `json:"query_filters"`
	Connector    Connector       `json:"connector"`
}

func NewDbQueryWhere(filters []DbQueryFilter, connector Connector) DbQueryWhere {
	return DbQueryWhere{
		QueryFilters: filters,
		Connector:    connector,
	}
}
