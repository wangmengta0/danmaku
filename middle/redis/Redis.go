package redis

import "github.com/redis/go-redis/v9"

var RdbReplay *redis.Client

func InitRedis() {
	RdbReplay = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}
