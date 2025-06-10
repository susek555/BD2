package notification

import (
	"fmt"

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
	dir := "DESC"
	if f.IsOrderDesc != nil && !*f.IsOrderDesc {
		dir = "ASC"
	}
	key := "notification_id"
	if f.OrderKey != nil {
		key = *f.OrderKey
	}

	if key == "notification_id" {
		query = query.Order(fmt.Sprintf("%s %s", key, dir))
	} else {
		query = query.Order(fmt.Sprintf("%s %s, notification_id %s", key, dir, dir))
	}

	if f.ReceiverID != nil {
		query = query.Where("user_id = ?", *f.ReceiverID)
	}
	if f.Seen != nil {
		query = query.Where("seen = ?", *f.Seen)
	}

	return query, nil
}
