package pagination

import (
	"gorm.io/gorm"
)

var (
	DEFAULT_PAGE_SIZE int = 8
	DEFAULT_PAGE      int = 1
)

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginationResponse struct {
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

func (pr *PaginationRequest) setDefaults(page int, pageSize int) {
	if pr.Page <= 0 {
		pr.Page = page
	}
	if pr.PageSize <= 0 {
		pr.PageSize = pageSize
	}

}

func Paginate(pr *PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pr.setDefaults(DEFAULT_PAGE, DEFAULT_PAGE_SIZE)
		offset := (pr.Page - 1) * pr.PageSize
		return db.Offset(offset).Limit(pr.PageSize)
	}
}
