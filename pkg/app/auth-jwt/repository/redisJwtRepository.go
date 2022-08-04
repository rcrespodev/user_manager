package repository

import (
	"encoding/json"
	"fmt"
	redisClient "github.com/go-redis/redis/v8"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/redis"
)

type RedisJwtRepository struct {
	redisRepository *redis.RedisRepository
}

func NewRedisJwtRepository(redisClient *redisClient.Client) *RedisJwtRepository {
	return &RedisJwtRepository{redisRepository: redis.NewRedisRepository(redisClient)}
}

func (r *RedisJwtRepository) Update(command jwtDomain.UpdateCommand, log *returnLog.ReturnLog) {
	key := r.buildKey(command.Command.Uuid)
	value, err := json.Marshal(command.Command)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	err = r.redisRepository.RedisCliente().Set(r.redisRepository.Ctx(), key, value, command.Command.Duration).Err()
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}
}

func (r *RedisJwtRepository) FindByUuid(query jwtDomain.FindByUuidQuery) *jwtDomain.JwtSchema {

	var jwtSchema *jwtDomain.JwtSchema
	key := r.buildKey(query.Uuid)
	result, err := r.redisRepository.RedisCliente().Get(r.redisRepository.Ctx(), key).Result()
	if err != nil {
		return nil
	}

	err = json.Unmarshal([]byte(result), &jwtSchema)
	if err != nil {
		return nil
	}

	return jwtSchema
}

func (r *RedisJwtRepository) buildKey(uuid string) string {
	return fmt.Sprintf("%s-%s", "jwt", uuid)
}
