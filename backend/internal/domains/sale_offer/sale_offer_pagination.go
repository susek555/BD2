package sale_offer

import paginator "github.com/pilagod/gorm-cursor-paginator"

type PagingQuery struct {
	After  *string
	Before *string
	Limit  *int
	Order  *string
}

func GetProductPaginator(q PagingQuery, orderKey string) *paginator.Paginator {
	p := paginator.New()

	p.SetKeys(orderKey, "ID")

	if q.After != nil {
		p.SetAfterCursor(*q.After)
	}

	if q.Before != nil {
		p.SetBeforeCursor(*q.Before)
	}

	if q.Limit != nil {
		p.SetLimit(*q.Limit)
	}

	if q.Order != nil && *q.Order == "asc" {
		p.SetOrder(paginator.ASC)
	}

	return p
}
