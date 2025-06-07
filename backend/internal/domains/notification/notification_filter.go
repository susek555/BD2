package notification

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type NotificationFilter struct {
	Pagination  pagination.PaginationRequest `json:"pagination"`
	OrderKey    *string                      `json:"order_key"`
	IsOrderDesc *bool                        `json:"is_order_desc,omitempty"`
	ReceiverID  *uint                        `json:"receiver_id,omitempty"`
	Seen        *bool                        `json:"seen,omitempty"`
}

func NewNotificationFilter() *NotificationFilter {
	return &NotificationFilter{}
}

func (f *NotificationFilter) ApplyNotificationFilters(query *gorm.DB) (*gorm.DB, error) {
	if f.OrderKey != nil {
		order := *f.OrderKey
		if f.IsOrderDesc != nil && *f.IsOrderDesc {
			order += " DESC"
		}
		query = query.Order(order)
	}

	if f.ReceiverID != nil {
		query = query.Where("user_id = ?", *f.ReceiverID)
	}
	if f.Seen != nil {
		query = query.Where("seen = ?", *f.Seen)
	}
	return query, nil
}
