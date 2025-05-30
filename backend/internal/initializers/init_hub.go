package initializers

import (
	"context"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
)

var Hub *ws.Hub

func InitializeHub() {
	Hub = ws.NewHub(ClientNotificationRepo, DB)
	go Hub.Run()
	ctx := context.Background()
	Hub.StartRedisFanIn(ctx, RedisClient)
}
