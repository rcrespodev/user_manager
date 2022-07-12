package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	userRepository "github.com/rcrespodev/user_manager/pkg/app/user/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command/factory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
)

var Instance *Kernel

type Kernel struct {
	commandBus        *command.Bus
	messageRepository message.MessageRepository
	userRepository    domain.UserRepository
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}
	Instance = &Kernel{
		messageRepository: repository.NewRedisMessageRepository(redisClient),
		userRepository:    userRepository.NewMySqlUserRepository(mySqlClient),
	}
	Instance.commandBus = factory.NewCommandBusInstance(factory.NewCommandBusCommand{
		UserRepository: Instance.userRepository,
	})
	return Instance
}

//func (k *Kernel) LoadCommandBus() {
//	k.commandBus = factory.NewCommandBusInstance()
//}

func (k Kernel) CommandBus() *command.Bus {
	return k.commandBus
}

//func (k Kernel) LoadMessageRepository() {
//	k.messageRepository = repository.NewRedisMessageRepository()
//}

func (k Kernel) MessageRepository() message.MessageRepository {
	return k.messageRepository
}

//func (k *Kernel) LoadUserRepository() {
//	k.userRepository = userRepository.NewMySqlUserRepository()
//}

func (k *Kernel) UserRepository() domain.UserRepository {
	return k.userRepository
}
