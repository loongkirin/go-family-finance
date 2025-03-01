package database

type DbQueryFilter struct {
	FieldName       string          `json:"field_name"`
	FilterValues    []interface{}   `json:"filter_values"`
	FilterOperation FilterOperation `json:"filter_operation"`
	Connector       Connector       `json:"connector"`
	FieldType       string          `json:"field_type"`
}

func NewDbQueryFilter(field string, values []interface{}, op FilterOperation, ft string) DbQueryFilter {
	return DbQueryFilter{
		FieldName:       field,
		FilterValues:    values,
		FilterOperation: op,
		Connector:       AND,
		FieldType:       ft,
	}
}
