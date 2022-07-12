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
	redisCliente *redis.Client
	ctx          context.Context
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	redisRepository := &RedisRepository{}
	switch redisClient {
	case nil:
		redisRepository.redisCliente = redisRepository.newConnection()
		redisRepository.ctx = context.Background()
	default:
		redisRepository.redisCliente = redisClient
		redisRepository.ctx = context.Background()
	}

	return redisRepository
}

func (r RedisRepository) RedisCliente() *redis.Client {
	return r.redisCliente
}

func (r RedisRepository) Ctx() context.Context {
	return r.ctx
}

func (r RedisRepository) newConnection() *redis.Client {
	redisConf := config.Conf.Redis
	redisCliente := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Host),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := redisCliente.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection %v", err)
	}

	return redisCliente
}
