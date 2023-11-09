package redis

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Redis struct is new api for go-redis library
type Redis struct {
	client *redis.Client
	ctx    context.Context
}

// CreateRedisClient return Redis structure with given options
func CreateRedisClient(Addr string, Password string, DB int) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     "redis-auth:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		ctx: context.Background(),
	}
}

// Set is function that set value with given key and expiration time
func (r *Redis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(r.ctx, key, value, expiration)
}

// Get is function that get value with given key
func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(r.ctx, key)
}

// Delete witll delete value with given key
func (r *Redis) Delete(key string) *redis.IntCmd {
	return r.client.Del(r.ctx, key)
}

// Contain will check if given key is present in DB
func (r *Redis) Contain(key string) (bool, error) {
	_, err := r.Get(key).Result()

	return err == nil, err
}
