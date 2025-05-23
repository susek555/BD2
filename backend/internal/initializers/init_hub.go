package initializers

import (
	"context"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
)

var Hub *auctionws.Hub

func InitializeHub() {
	Hub = auctionws.NewHub()
	go Hub.Run()
	ctx := context.Background()
	Hub.StartRedisFanIn(ctx, RedisClient)
}
