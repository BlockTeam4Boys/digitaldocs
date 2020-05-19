package repository

import "fmt"

func userKeyKey(userID uint) string {
	const pattern = "user_id:%v:key"
	return fmt.Sprintf(pattern, userID)
}

func (r *RedisTokenRepository) SetUserKey(userID uint, key []byte) error {
	return r.rdb.Set(userKeyKey(userID), key, 0).Err()
}

func (r *RedisTokenRepository) GetUserKey(userID uint) (string, error) {
	return r.rdb.Get(userKeyKey(userID)).Result()
}
