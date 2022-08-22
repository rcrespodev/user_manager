package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"log"
	_ "os"
)

type RedisRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	redisRepository := &RedisRepository{}
	switch redisClient {
	case nil:
		redisRepository.redisClient = redisRepository.newConnection()
		redisRepository.ctx = context.Background()
	default:
		redisRepository.redisClient = redisClient
		redisRepository.ctx = context.Background()
	}

	return redisRepository
}

func (r RedisRepository) RedisClient() *redis.Client {
	return r.redisClient
}

func (r RedisRepository) Ctx() context.Context {
	return r.ctx
}

func (r RedisRepository) newConnection() *redis.Client {
	redisConf := config.Conf.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Host),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection %v", err)
	}

	return redisClient
}
