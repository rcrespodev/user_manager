package userSessionRepository

import (
	"encoding/json"
	"fmt"
	redisClient "github.com/go-redis/redis/v8"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/redis"
)

const (
	userSessionKey = "userSession"
)

type RedisUserSessionRepository struct {
	redisRepository *redis.RedisRepository
}

func NewRedisUserSessionRepository(redisClient *redisClient.Client) *RedisUserSessionRepository {
	return &RedisUserSessionRepository{
		redisRepository: redis.NewRedisRepository(redisClient),
	}
}

func (r *RedisUserSessionRepository) UpdateUserSession(command domain.UpdateUserSessionCommand, log *returnLog.ReturnLog) {
	key := r.buildKey(command.UserUuid)
	ctx := r.redisRepository.Ctx()

	userSessionSchema := domain.UserSessionSchema{
		UserUuid:     command.UserUuid,
		IsLogged:     command.IsLogged,
		LastLoginOn:  command.LastLoginOn,
		LastLogoutOn: command.LastLogoutOn,
	}
	valueBytes, err := json.Marshal(userSessionSchema)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	err = r.redisRepository.RedisCliente().Set(ctx, key, valueBytes, 0).Err()
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}
}

func (r RedisUserSessionRepository) GetUserSession(userUuid string) *domain.UserSessionSchema {
	var userSessionSchema *domain.UserSessionSchema
	key := r.buildKey(userUuid)
	ctx := r.redisRepository.Ctx()
	result, err := r.redisRepository.RedisCliente().Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	err = json.Unmarshal([]byte(result), &userSessionSchema)
	if err != nil {
		return nil
	}

	return userSessionSchema
}

func (r RedisUserSessionRepository) buildKey(userUuid string) string {
	return fmt.Sprintf("%s-%s", userSessionKey, userUuid)
}
