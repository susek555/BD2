package ws

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func PublishAuctionEvent(
	ctx context.Context,
	rdb *redis.Client,
	offerID string,
	envelope *Envelope,
) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	return rdb.Publish(ctx, "offer."+offerID, data).Err()
}
