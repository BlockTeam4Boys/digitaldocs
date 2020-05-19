package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisTokenRepository struct {
	rdb              *redis.Client
	refreshTokenTime time.Duration
	accessTokenTime  time.Duration
}

func NewRedisSessionRepository(rdb *redis.Client, accessTokenTime, refreshTokenTime time.Duration) *RedisTokenRepository {
	return &RedisTokenRepository{
		rdb:              rdb,
		accessTokenTime:  accessTokenTime,
		refreshTokenTime: refreshTokenTime,
	}
}
