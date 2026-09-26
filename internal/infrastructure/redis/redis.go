package redis

import (
	"context"
	"log"

	redis "github.com/redis/go-redis/v9"
)

var Client *redis.Client

// Connect initializes the Redis client and pings it to ensure connection is alive.
func Connect(redisURL string) {
	Client = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set by default
		DB:       0,  // use default DB
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisURL, err)
	}

	log.Println("Connected to Redis successfully")
}
