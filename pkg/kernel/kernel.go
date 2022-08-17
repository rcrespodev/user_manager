package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command/commandBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query/queryBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/log/file"
)

var Instance *Kernel

type Kernel struct {
	commandBus        *command.Bus
	queryBus          *query.Bus
	messageRepository message.MessageRepository
	userRepository    domain.UserRepository
	config            *config.Config
	jwt               *jwtDomain.Jwt
	jwtRepository     jwtDomain.JwtRepository
	logFile           *file.LogFile
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}

	// log File instance
	Instance.logFile = file.NewLogFile(config.Conf.Log.File.Path)

	// jwt instance
	jwt, jwtRepository := jwtFactory(redisClient)
	Instance.jwt = jwt
	Instance.jwtRepository = jwtRepository

	// repositories instance
	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)

	// command bus instance
	Instance.commandBus = commandBusFactory.NewCommandBusInstance(commandBusFactory.NewCommandBusCommand{
		UserRepository: Instance.userRepository,
		Jwt:            Instance.jwt,
		JwtRepository:  Instance.jwtRepository,
	})

	// query bus instance
	Instance.queryBus = queryBusFactory.NewQueryBusInstance(
		queryBusFactory.NewQueryBusCommand{
			UserRepository: Instance.userRepository})

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

func (k *Kernel) QueryBus() *query.Bus {
	return k.queryBus
}
