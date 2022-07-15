package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	userRepository "github.com/rcrespodev/user_manager/pkg/app/user/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
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
	config            *config.Config
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}
	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)
	Instance.commandBus = factory.NewCommandBusInstance(factory.NewCommandBusCommand{
		UserRepository: Instance.userRepository,
	})
	return Instance
}

func (k Kernel) CommandBus() *command.Bus {
	return k.commandBus
}

func (k Kernel) MessageRepository() message.MessageRepository {
	return k.messageRepository
}

func (k *Kernel) UserRepository() domain.UserRepository {
	return k.userRepository
}

func (k Kernel) Config() *config.Config {
	return k.config
}
