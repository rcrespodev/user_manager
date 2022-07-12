package repository

import (
	"encoding/json"
	"fmt"
	redisClient "github.com/go-redis/redis/v8"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/redis"
	"log"
)

type RedisMessageRepository struct {
	redisRepository *redis.RedisRepository
}

func NewRedisMessageRepository(redisClient *redisClient.Client) *RedisMessageRepository {
	messages := GetMessagesSourceData()

	r := &RedisMessageRepository{
		redisRepository: redis.NewRedisRepository(redisClient),
	}

	for _, messageSchema := range messages.Messages {
		ctx := r.redisRepository.Ctx()
		key := r.buildKey(messageSchema.Id, messageSchema.Pkg)
		messageBytes, err := json.Marshal(messageSchema)
		if err != nil {
			log.Fatal(err)
		}

		err = r.redisRepository.RedisCliente().Set(ctx, key, messageBytes, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}

func (r RedisMessageRepository) GetMessageData(id message.MessageId, messagePkg string) (text string, clientErrorType message.ClientErrorType) {
	var messageSchema MessageSchema
	key := r.buildKey(id, messagePkg)
	ctx := r.redisRepository.Ctx()
	result, err := r.redisRepository.RedisCliente().Get(ctx, key).Result()
	err = json.Unmarshal([]byte(result), &messageSchema)
	if err != nil {
		return "", 0
	}

	return messageSchema.Text, messageSchema.ClientErrorType
}

func (r RedisMessageRepository) buildKey(id message.MessageId, pkg string) string {
	return fmt.Sprintf("%v-%v", id, pkg)
}
