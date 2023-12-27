package cashredis

import (
	"tg-bot-golang/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewUniversalRedisClient(cfg *config.ConfigRedis) redis.UniversalClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB, // use default DB
	})

	return rdb
}
