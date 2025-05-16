package review

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type ReviewFilter struct {
	Pagination  pagination.PaginationRequest `json:"pagination"`
	OrderKey    *string                      `json:"order_key"`
	IsOrderDesc *bool                        `json:"is_order_desc"`
	Ratings     *[]uint                       `json:"ratings"`
}

func NewReviewFilter() *ReviewFilter {
	return &ReviewFilter{
		Ratings: &[]uint{},
	}
}

func (f *ReviewFilter) ApplyReviewFilters(query *gorm.DB) (*gorm.DB, error) {
	if f.OrderKey != nil {
		order := *f.OrderKey
		if f.IsOrderDesc != nil && *f.IsOrderDesc {
			order += " DESC"
		}
		query = query.Order(order)
	}

	if len(*f.Ratings) > 0 {
		query = query.Where("rating IN ?", *f.Ratings)
	}

	return query, nil
}


