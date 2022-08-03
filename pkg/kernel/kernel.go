package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
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
	jwt               *jwtDomain.Jwt
	jwtRepository     jwtDomain.JwtRepository
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}

	jwt, jwtRepository := jwtFactory(redisClient)
	Instance.jwt = jwt
	Instance.jwtRepository = jwtRepository

	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)
	Instance.commandBus = factory.NewCommandBusInstance(factory.NewCommandBusCommand{
		UserRepository: Instance.userRepository,
		Jwt:            Instance.jwt,
		JwtRepository:  Instance.jwtRepository,
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

func (k *Kernel) Jwt() *jwtDomain.Jwt {
	return k.jwt
}

func (k *Kernel) JwtRepository() jwtDomain.JwtRepository {
	return k.jwtRepository
}
