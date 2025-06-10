package initializers

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis(ctx context.Context) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}
	log.Println("Connected to Redis")
	RedisClient = redisClient
}
