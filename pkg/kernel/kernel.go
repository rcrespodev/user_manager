package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/infrastructure"
	emailSenderRepository "github.com/rcrespodev/user_manager/pkg/app/emailSender/repository"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/amqp/rabbitMq"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command/commandBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event/consumersFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event/eventBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query/queryBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/jwt"
	"github.com/rcrespodev/user_manager/pkg/kernel/log/file"
)

var Instance *Kernel

// Kernel implements a singleton pattern, and contains all dependencies of application.
// Also, Kernel has the responsibility of instance the next deps:
// - cqrs: CommandBus, EventBus and QueryBus
// - Db: Mysql and Redis repositories.
// - RabbitMQ
// - Config
// - Jwt services with private and public rsa key.
// - EmailSender interface.
// - LogFile services.
// Kernel just receive DB clients. If any client == nil, the connection is instantiated
// by the concrete type of implementation e.g. MySqlRepository, RedisRepository.
type Kernel struct {
	commandBus          *command.AppBus
	queryBus            *query.Bus
	eventBus            event.Bus
	rabbitClient        *rabbitMq.Client
	messageRepository   message.MessageRepository
	userRepository      domain.UserRepository
	config              *config.Config
	jwt                 *jwtDomain.Jwt
	jwtRepository       jwtDomain.JwtRepository
	emailSender         emailSenderDomain.EmailSender
	sentEmailRepository emailSenderDomain.SentEmailRepository
	logFile             *file.LogFile
}

func NewPrdKernel(mySqlClient *sql.DB, redisClient *redis.Client, rabbitMqConnection *amqp.Connection) *Kernel {
	if Instance != nil {
		return Instance
	}

	Instance = &Kernel{
		config: config.Setup(),
	}

	// log File instance
	Instance.newLogFileService()

	// repositories instance
	Instance.newRepositories(mySqlClient, redisClient)

	// jwt instance
	Instance.newJwt(redisClient)

	// email sender dependencies
	Instance.newEmailSender(mySqlClient)

	// instance rabbitmq
	Instance.newRabbit(rabbitMqConnection)

	// buses instance
	Instance.newBuses(rabbitMqConnection)

	// Subscribe consumers to event bus
	consumersFactory.SubscribeConsumers(consumersFactory.SubscribeConsumersCommand{
		EventBus:            Instance.eventBus,
		CommandBus:          Instance.commandBus,
		UserRepository:      Instance.userRepository,
		EmailSender:         Instance.emailSender,
		SentEmailRepository: Instance.sentEmailRepository,
		MessageRepository:   Instance.messageRepository,
		WelcomeTemplatePath: config.Conf.Smtp.Welcome.Template,
	})

	return Instance
}

func (k *Kernel) newRepositories(mySqlClient *sql.DB, redisClient *redis.Client) {
	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)
	return
}

func (k *Kernel) newEmailSender(mySqlClient *sql.DB) {
	Instance.emailSender = infrastructure.NewSmtpEmailSender(infrastructure.EmailAuthConf{
		Host:     config.Conf.Smtp.Host,
		From:     config.Conf.Smtp.Username,
		Password: config.Conf.Smtp.Password,
		Port:     config.Conf.Smtp.Port,
	})
	Instance.sentEmailRepository = emailSenderRepository.NewMySqlSentEmailRepository(mySqlClient)
	return
}

func (k *Kernel) newJwt(redisClient *redis.Client) {
	j, jwtRepository := jwt.Factory(redisClient)
	Instance.jwt = j
	Instance.jwtRepository = jwtRepository
	return
}

func (k *Kernel) newRabbit(rabbitMqConnection *amqp.Connection) {
	Instance.rabbitClient = rabbitMq.NewRabbitMqClient(rabbitMqConnection)
	return
}

func (k *Kernel) newBuses(rabbitMqConnection *amqp.Connection) {
	// event bus instance
	Instance.eventBus = eventBusFactory.NewEventBusInstance(eventBusFactory.NewEventBusCommand{
		RabbitMqConnection: rabbitMqConnection})

	// command bus instance
	Instance.commandBus = commandBusFactory.NewCommandBusInstance(commandBusFactory.NewCommandBusCommand{
		User: struct {
			UserRepository domain.UserRepository
		}{UserRepository: Instance.userRepository},
		Jwt: struct {
			Jwt           *jwtDomain.Jwt
			JwtRepository jwtDomain.JwtRepository
		}{
			Jwt:           Instance.jwt,
			JwtRepository: Instance.jwtRepository,
		},
		EmailSender: struct {
			EmailSender         emailSenderDomain.EmailSender
			SentEmailRepository emailSenderDomain.SentEmailRepository
			WelcomeTemplatePath string
		}{
			EmailSender:         Instance.emailSender,
			SentEmailRepository: Instance.sentEmailRepository,
			WelcomeTemplatePath: config.Conf.Smtp.Welcome.Template,
		},
		EventBus: Instance.eventBus,
	})

	// query bus instance
	Instance.queryBus = queryBusFactory.NewQueryBusInstance(
		queryBusFactory.NewQueryBusCommand{
			UserRepository: Instance.userRepository})

	return
}

func (k *Kernel) newLogFileService() {
	Instance.logFile = file.NewLogFile(config.Conf.Log.File.Path)
	return
}

func (k *Kernel) CommandBus() *command.AppBus {
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

func (k *Kernel) RabbitClient() *rabbitMq.Client {
	return k.rabbitClient
}

func (k *Kernel) EventBus() event.Bus {
	return k.eventBus
}

func (k *Kernel) EmailSender() emailSenderDomain.EmailSender {
	return k.emailSender
}

func (k *Kernel) SentEmailRepository() emailSenderDomain.SentEmailRepository {
	return k.sentEmailRepository
}
