package pkg

import (
	"context"
	"fmt"
	"github.com/anhsbolic/closure-table-go/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	// Get Config
	env := config.GetEnvConfig()

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.Get("REDIS_HOST"), env.Get("REDIS_PORT")),
		Password: env.Get("REDIS_PASSWORD").(string),
		DB:       0,
	})

	// Check Connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		PanicIfError(err)
	}

	return client
}
