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
	"github.com/vrischmann/envconfig"
	"log"
)

var Instance *Kernel

type Kernel struct {
	commandBus        *command.Bus
	messageRepository message.MessageRepository
	userRepository    domain.UserRepository
	config            *Config
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	var config *Config
	if err := envconfig.Init(&config); err != nil {
		log.Fatal(err)
	}

	Instance = &Kernel{
		messageRepository: repository.NewRedisMessageRepository(redisClient),
		userRepository:    userRepository.NewMySqlUserRepository(mySqlClient),
		config:            config,
	}
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

func (k Kernel) Config() *Config {
	return k.config
}
