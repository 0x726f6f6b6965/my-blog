package client

import (
	"context"
	"fmt"

	"github.com/0x726f6f6b6965/my-blog/lib/config"
	"github.com/redis/go-redis/v9"
)

const (
	UserToken string = "email:%s"
	Secret    string = "secret"
)

func InitRedisClient(cfg *config.RedisCfg) *redis.Client {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	if cfg.MaxRetries != 0 {
		opt.MaxRetries = cfg.MaxRetries
	}
	return redis.NewClient(opt)
}

func GetToken(ctx context.Context, email string, rds *redis.Client) (string, error) {
	return rds.Get(ctx, fmt.Sprintf(UserToken, email)).Result()
}

func GetSecret(ctx context.Context, rds *redis.Client) (string, error) {
	return rds.Get(ctx, Secret).Result()
}
