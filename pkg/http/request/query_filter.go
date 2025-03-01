package request

type QueryFilter struct {
	FieldName    string        `json:"field_name"`
	FilterValues []interface{} `json:"filter_values"`
	Operator     Operator      `json:"operator"`
	Connector    Connector     `json:"connector"`
}

func NewQueryFilter(field string, values []interface{}, op Operator) *QueryFilter {
	return &QueryFilter{
		FieldName:    field,
		FilterValues: values,
		Operator:     op,
		Connector:    AND,
	}
}

func NewQueryFilterWithConnector(field string, values []interface{}, op Operator, connector Connector) *QueryFilter {
	return &QueryFilter{
		FieldName:    field,
		FilterValues: values,
		Operator:     op,
		Connector:    connector,
	}
}
