package request

type Operator string

const (
	EQ      Operator = "EQ"
	NEQ     Operator = "NEQ"
	LT      Operator = "LT"
	LTE     Operator = "LTE"
	GT      Operator = "GT"
	GTE     Operator = "GTE"
	LIKE    Operator = "LIKE"
	IN      Operator = "IN"
	BETWEEN Operator = "BETWEEN"
)

type Connector string

const (
	AND Connector = "AND"
	OR  Connector = "OR"
)

type Query struct {
	QueryWheres []*QueryWhere   `json:"query_wheres"`
	OrderBy     []*QueryOrderBy `json:"order_by"`
	PageSize    int             `json:"page_size"`
	PageNumber  int             `json:"page_number"`
}

func NewQuery(wheres []*QueryWhere, ps int, pn int, order []*QueryOrderBy) *Query {
	return &Query{
		QueryWheres: wheres,
		PageSize:    ps,
		PageNumber:  pn,
		OrderBy:     order,
	}
}
