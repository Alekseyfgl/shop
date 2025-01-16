package model

type Paginate[T any] struct {
	PageNumber     int `json:"pageNumber"`
	RowTotalCount  int `json:"rowTotalCount"`
	TotalPageCount int `json:"totalPageCount"`
	PageSize       int `json:"pageSize"`
	Items          []T `json:"items"`
}
