package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type RedisRepository struct {
	redisCliente *redis.Client
	ctx          context.Context
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	switch redisClient {
	case nil:
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		if host == "" || port == "" {
			log.Fatalf("redis env vars nots found")
		}

		r := &RedisRepository{
			redisCliente: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%v:%v", host, port),
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
			ctx: context.Background(),
		}
		if err := r.redisCliente.Ping(r.ctx).Err(); err != nil {
			log.Fatalf("Redis %v", err)
		}
		return r
	default:
		return &RedisRepository{
			redisCliente: redisClient,
			ctx:          context.Background(),
		}
	}

}

func (r RedisRepository) RedisCliente() *redis.Client {
	return r.redisCliente
}

func (r RedisRepository) Ctx() context.Context {
	return r.ctx
}
