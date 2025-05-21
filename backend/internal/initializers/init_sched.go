package initializers

import (
	"context"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
)

var Sched scheduler.SchedulerInterface

func InitializeScheduler() {
	repo := bid.NewBidRepository(DB)
	notificationRepo := notification.NewNotificationRepository(DB)
	saleOfferRepo := sale_offer.NewSaleOfferRepository(DB)
	manufacturerRepo := manufacturer.NewManufacturerRepository(DB)
	manufacturerService := manufacturer.NewManufacturerService(manufacturerRepo)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerService)
	notificationService := notification.NewNotificationService(notificationRepo, saleOfferService)
	Sched = scheduler.NewScheduler(repo, RedisClient, notificationService)
	go Sched.Run(context.Background())
}
