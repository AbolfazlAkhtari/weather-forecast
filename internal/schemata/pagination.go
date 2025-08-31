package schemata

const PaginationLimit = 10

type Pagination struct {
	TotalPage   int64 `json:"total_page"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int   `json:"current_page"`
}
