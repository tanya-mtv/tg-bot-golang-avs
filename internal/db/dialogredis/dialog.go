package dialogredis

import (
	"context"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type DPoint struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewDialogStorage(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *DPoint {
	return &DPoint{log: log, cfg: cfg, redisClient: redisClient}
}

func (d *DPoint) Add(ctx context.Context, user, value string) {
	if err := d.redisClient.Set(ctx, user, value, time.Duration(24)*time.Hour).Err(); err != nil {
		d.log.Warnf("Can't set dialog point for user", err)
	}
	d.log.Infof("Set dialog point for user: %s %s", user, value)
}

func (d *DPoint) Del(ctx context.Context, user string) {
	if err := d.redisClient.Del(ctx, user).Err(); err != nil {
		d.log.Warnf("Can't delete dialog point for user", err)
	}
	d.log.Infof("Delete dialog point for user: %s %s", user)
}

func (d *DPoint) Get(ctx context.Context, user string) string {

	point, err := d.redisClient.Get(ctx, user).Bytes()

	if err != nil {
		if err != redis.Nil {
			d.log.WarnMsg("redisClient.Get. Error getting user dialogs point", err)
			return ""
		}
		return ""
	}
	d.log.Infof("Get dialog point for user: %s %s", user, point)
	return string(point)

}
