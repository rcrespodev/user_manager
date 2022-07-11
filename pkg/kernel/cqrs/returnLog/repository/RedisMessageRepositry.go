package repository

import (
	"encoding/json"
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository"
	"log"
	"os"
)

type RedisMessageRepository struct {
	redisRepository *repository.RedisRepository
}

type MessagesSchema struct {
	Messages []MessageSchema `json:"messages"`
}

type MessageSchema struct {
	Id              message.MessageId       `json:"id"`
	Pkg             string                  `json:"pkg"`
	Text            string                  `json:"text"`
	ClientErrorType message.ClientErrorType `json:"client_error_type"`
}

func NewRedisMessageRepository() *RedisMessageRepository {
	sourcePath := os.Getenv("JSON_MESSAGES")
	if sourcePath == "" {
		log.Fatal("env JSON_MESSAGES not found")
	}

	//sourceData, err := ioutil.ReadFile("messages.json")
	//if err != nil {
	//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf(dir)
	//	pc, file, line, ok := runtime.Caller(0)
	//	log.Printf("%v, %v, %v, %v", pc, file, line, ok)
	//	log.Fatal(err)
	//}
	//
	//var messages MessagesSchema
	//if err = json.Unmarshal(sourceData, &messages); err != nil {
	//	log.Fatal(err)
	//}
	//
	r := &RedisMessageRepository{
		redisRepository: repository.NewRedisRepository(),
	}

	//for _, messageSchema := range messages.Messages {
	//	ctx := r.redisRepository.Ctx()
	//	key := r.buildKey(messageSchema.Id, messageSchema.Pkg)
	//	messageBytes, err := json.Marshal(messageSchema)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	err = r.redisRepository.RedisCliente().Set(ctx, key, messageBytes, 0).Err()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
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
