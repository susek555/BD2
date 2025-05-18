package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
)

var Sched scheduler.SchedulerInterface

func InitializeScheduler() {
	repo := bid.NewBidRepository(DB)
	Sched = scheduler.NewScheduler(DB, repo, RedisClient)
}
