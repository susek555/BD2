package pagination

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrPageOutOfRange error = errors.New("page out of range")
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
	TotalPages   int64 `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

func (pr *PaginationRequest) setDefaults(page int, pageSize int) {
	if pr.Page <= 0 {
		pr.Page = page
	}
	if pr.PageSize <= 0 {
		pr.PageSize = pageSize
	}
}

func (pr *PaginationRequest) CalculateTotalPages(totalRecords int64) int64 {
	pr.setDefaults(DEFAULT_PAGE, DEFAULT_PAGE_SIZE)
	totalPages := totalRecords / int64(pr.PageSize)
	if totalPages%int64(pr.PageSize) != 0 {
		totalPages++
	}
	return totalPages
}

func (pr *PaginationRequest) ValidatePageNumber(totalPages int64) error {
	if pr.Page > int(totalPages) {
		return ErrPageOutOfRange
	}
	return nil
}

func Paginate(pr *PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pr.Page - 1) * pr.PageSize
		return db.Offset(offset).Limit(pr.PageSize)
	}
}
