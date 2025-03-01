package request

type GetDataByQueryRequest struct {
	Query *Query `json:"query"`
}

type DataRequest[T any] struct {
	Data T `json:"data"`
}

type DataListRequest[T any] struct {
	DataList []T `json:"data_list"`
}
