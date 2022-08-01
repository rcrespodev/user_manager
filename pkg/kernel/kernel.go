package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
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
	jwtConfig         *jwtDomain.JwtConfig
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}
	Instance.jwtConfig = jwtDomain.NewJwtConfig(config.Conf.Jwt.Secret, config.Conf.Jwt.ExpirationTime)
	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)
	Instance.commandBus = factory.NewCommandBusInstance(factory.NewCommandBusCommand{
		UserRepository: Instance.userRepository,
	})
	return Instance
}

func (k *Kernel) CommandBus() *command.Bus {
	return k.commandBus
}

func (k *Kernel) MessageRepository() message.MessageRepository {
	return k.messageRepository
}

func (k *Kernel) UserRepository() domain.UserRepository {
	return k.userRepository
}

func (k *Kernel) Config() *config.Config {
	return k.config
}

func (k *Kernel) JwtConfig() *jwtDomain.JwtConfig {
	return k.jwtConfig
}
