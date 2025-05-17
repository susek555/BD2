package auctionws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func publishAuctionEvent(
	ctx context.Context,
	rdb *redis.Client,
	auctionID string,
	envelope *Envelope,
) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		log.Println("Error marshalling envelope:", err)
	}
	return rdb.Publish(ctx, "auction."+auctionID, data).Err()
}
