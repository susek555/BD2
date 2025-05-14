package pagination

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrPageOutOfRange   = errors.New("page out of range")
	ErrNegativePageSize = errors.New("page size must be greater than zero")
)

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginationResponse struct {
	TotalPages   int64 `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

func (pr *PaginationRequest) validatePageNumber(totalPages int64) error {
	if pr.Page < 1 || int64(pr.Page) > totalPages {
		return ErrPageOutOfRange
	}
	return nil
}

func (pr *PaginationRequest) validatePageSize() error {
	if pr.PageSize <= 0 {
		return ErrNegativePageSize
	}
	return nil
}

func (pr *PaginationRequest) calculateTotalPages(totalRecords int64) int64 {
	totalPages := totalRecords / int64(pr.PageSize)
	if totalRecords%int64(pr.PageSize) != 0 || totalRecords == 0 {
		totalPages++
	}
	return totalPages
}

func Paginate(pr *PaginationRequest, totalRecords int64) (func(db *gorm.DB) *gorm.DB, *PaginationResponse, error) {
	if err := pr.validatePageSize(); err != nil {
		return nil, nil, err
	}
	totalPages := pr.calculateTotalPages(totalRecords)
	if err := pr.validatePageNumber(totalPages); err != nil {
		return nil, nil, err
	}
	offset := (pr.Page - 1) * pr.PageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pr.PageSize)
	}, &PaginationResponse{TotalPages: totalPages, TotalRecords: totalRecords}, nil
}
