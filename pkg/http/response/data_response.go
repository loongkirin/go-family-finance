package response

type Pagination struct {
	PageSize    int  `json:"page_size"`
	PageNumber  int  `json:"page_number"`
	HasNextPage bool `json:"has_next_page"`
}

type DataResponse[T any] struct {
	Data T `json:"data"`
}

type DataListResponse[T any] struct {
	DataList   []T        `json:"data_list"`
	Pagination Pagination `json:"page_info"`
}
