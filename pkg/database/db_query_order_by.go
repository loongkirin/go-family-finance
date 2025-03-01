package database

type DbQueryOrderBy struct {
	FieldName string `json:"field_name"`
	IsAsc     bool   `json:"is_asc"`
}

func NewDDbQueryOrderBy(field string, asc bool) DbQueryOrderBy {
	return DbQueryOrderBy{
		FieldName: field,
		IsAsc:     asc,
	}
}
