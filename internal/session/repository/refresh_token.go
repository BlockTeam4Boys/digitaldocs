package repository

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	TokensNotEqualErr = fmt.Errorf("refresh and current tokens are not equal to the token in db")
	TokenExpiredErr   = fmt.Errorf("token expired in db")
)

func (r *RedisTokenRepository) SetRefreshTokenHash(sessionID, token string) error {
	return r.rdb.Set(sessionID, token, r.refreshTokenTime).Err()
}

func (r *RedisTokenRepository) GetRefreshTokenHash(sessionID string) (string, error) {
	return r.rdb.Get(sessionID).Result()
}

func (r *RedisTokenRepository) DeleteRefreshToken(sessionID string) error {
	return r.rdb.Del(sessionID).Err()
}

func (r *RedisTokenRepository) CheckAndSetRefreshToken(sessionID, parentHash, currentHash string) error {
	txf := func(tx *redis.Tx) error {
		dbHash, err := tx.Get(sessionID).Result()
		if err != nil {
			return err
		}
		if parentHash == dbHash {
			_, err := tx.TxPipelined(func(pipe redis.Pipeliner) error {
				pipe.Set(sessionID, currentHash, r.refreshTokenTime)
				return nil
			})
			return err
		}
		if currentHash == dbHash {
			return nil
		}
		return TokensNotEqualErr
	}
	err := r.rdb.Watch(txf, sessionID)
	return err
}
