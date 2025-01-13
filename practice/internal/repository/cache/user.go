package cache

import (
	"basic-go/practice/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, id int64, user domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration //过期时间，也可以在Get中传入
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15, //这里写死了，也可以传入
	}
}

// 如果err为nil就认为缓存中有数据
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, id int64, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(id)
	err = cache.client.Set(ctx, key, string(val), cache.expiration).Err()
	return err
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
