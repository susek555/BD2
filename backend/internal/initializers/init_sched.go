package initializers

import (
	"context"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
)

var Sched scheduler.SchedulerInterface

func InitializeScheduler() {
	Sched = scheduler.NewScheduler(BidRepo, RedisClient, NotificationService, SaleOfferRepo, Hub)
	if err := Sched.LoadAuctions(); err != nil {
		panic("failed to load auctions: " + err.Error())
	}
	go Sched.Run(context.Background())
}
