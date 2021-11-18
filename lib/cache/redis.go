package redis

import (
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
)


var redisClient *redis.Client
var redisConnectOnce sync.Once


func GetRedisClient() *redis.Client {
	redisConnectOnce.Do(func() {
		println("Calling Redis .Once ......", os.Getenv("REDIS_HOST") != "")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"), // no password set
			DB:       0,  // use default DB
		})
		println("client", redisClient)
	})
	return redisClient
}
