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
	"log"
	"os"
	"strconv"
	"time"
)

var Instance *Kernel

type Kernel struct {
	commandBus        *command.Bus
	messageRepository message.MessageRepository
	userRepository    domain.UserRepository
	config            *config.Config
	jwt               *jwtDomain.Jwt
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}
	certPublicKey, err := os.ReadFile(config.Conf.Jwt.Key.Public)
	if err != nil {
		log.Fatal(err)
	}
	certPrivateKey, err := os.ReadFile(config.Conf.Jwt.Key.Private)
	if err != nil {
		log.Fatal(err)
	}
	expirationTime := config.Conf.Jwt.ExpirationTime
	expirationTimeInt, err := strconv.Atoi(expirationTime)
	if err != nil {
		log.Fatal(err)
	}
	Instance.jwt = jwtDomain.NewJwt(certPublicKey, certPrivateKey, time.Duration(expirationTimeInt))

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

func (k *Kernel) Jwt() *jwtDomain.Jwt {
	return k.jwt
}
