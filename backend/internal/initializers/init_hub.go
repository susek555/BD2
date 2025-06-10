package initializers

import (
	"context"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
)

var Hub ws.HubInterface

func InitializeHub() {
	Hub = ws.NewHub(NotificationService, UserOfferRepo)
	go Hub.Run()
	ctx := context.Background()
	Hub.StartRedisFanIn(ctx, RedisClient)
}
