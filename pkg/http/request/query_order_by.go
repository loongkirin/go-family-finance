package request

type QueryOrderBy struct {
	FieldName string `json:"field_name"`
	IsAsc     bool   `json:"is_asc"`
}

func NewQueryOrderBy(field string, asc bool) *QueryOrderBy {
	return &QueryOrderBy{
		FieldName: field,
		IsAsc:     asc,
	}
}
