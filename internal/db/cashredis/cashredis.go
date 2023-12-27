package cashredis

import (
	"context"
	"encoding/json"
	"fmt"
	"tg-bot-golang/internal/appmodels.go"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisStorage(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *RedisStorage {
	return &RedisStorage{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *RedisStorage) PutProduct(ctx context.Context, key string, product *appmodels.Product) {
	fmt.Printf("Product %+v \n", product)
	productBytes, err := json.Marshal(product)
	if err != nil {
		r.log.Infof("json.Marshal", err)
		return
	}

	if err := r.redisClient.Set(ctx, key, productBytes, time.Duration(24)*time.Hour).Err(); err != nil {
		r.log.Warnf("Can't set value to redis", err)
	}
	r.log.Infof("Set value in redis with key: %s", key)
}

func (r *RedisStorage) GetProduct(ctx context.Context, key string) (*appmodels.Product, error) {
	productBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.Get", err)
		}
		return nil, err
	}

	var product appmodels.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, err
	}

	r.log.Debugf("Get product from redis, key: %s", key)
	return &product, nil
}
