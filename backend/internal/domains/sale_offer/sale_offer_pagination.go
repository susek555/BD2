package sale_offer

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

type PagingQuery struct {
	Cursor paginator.Cursor `json:"cursor"`
	Order  *paginator.Order `json:"order"`
	Limit  *int             `json:"limit"`
}

var OrderMap = map[string]string{
	"Price":       "price",
	"DateOfIssue": "date_of_issue",
}

func GetOfferPaginator(q PagingQuery, orderKey *string) *paginator.Paginator {
	var rules []paginator.Rule
	if orderKey == nil {
		rules = []paginator.Rule{{Key: "ID"}}
	} else {
		rules = []paginator.Rule{{Key: *orderKey, SQLRepr: OrderMap[*orderKey]}, {Key: "ID"}}
	}
	cfg := paginator.Config{
		Rules:         rules,
		Limit:         8,
		Order:         paginator.DESC,
		AllowTupleCmp: paginator.TRUE,
	}
	p := paginator.New(&cfg)

	if q.Cursor.After != nil {
		p.SetAfterCursor(*q.Cursor.After)
	}

	if q.Cursor.Before != nil {
		p.SetBeforeCursor(*q.Cursor.Before)
	}

	if q.Limit != nil {
		p.SetLimit(*q.Limit)
	}

	if q.Order != nil {
		p.SetOrder(*q.Order)
	}
	return p
}
