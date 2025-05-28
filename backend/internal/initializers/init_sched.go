package initializers

import (
	"context"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
)

var Sched scheduler.SchedulerInterface

func InitializeScheduler() {
	Sched = scheduler.NewScheduler(BidRepo, RedisClient, NotificationService, SaleOfferRepo, Hub)
	go Sched.Run(context.Background())
}
