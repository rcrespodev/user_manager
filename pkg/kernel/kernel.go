package kernel

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/events"
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
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query/queryBusFactory"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/log/file"
	"log"
)

var Instance *Kernel

type Kernel struct {
	commandBus          *command.Bus
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
	Instance.logFile = file.NewLogFile(config.Conf.Log.File.Path)

	// jwt instance
	jwt, jwtRepository := jwtFactory(redisClient)
	Instance.jwt = jwt
	Instance.jwtRepository = jwtRepository

	// repositories instance
	Instance.messageRepository = repository.NewRedisMessageRepository(redisClient)
	Instance.userRepository = userRepository.NewMySqlUserRepository(mySqlClient)

	// email sender dependencies
	Instance.emailSender = infrastructure.NewSmtpEmailSender(infrastructure.EmailAuthConf{
		Host:     config.Conf.Smtp.Host,
		From:     config.Conf.Smtp.Username,
		Password: config.Conf.Smtp.Password,
		Port:     config.Conf.Smtp.Port,
	})
	Instance.sentEmailRepository = emailSenderRepository.NewMySqlSentEmailRepository(mySqlClient)

	// event bus instance
	Instance.rabbitClient = rabbitMq.NewRabbitMqClient(rabbitMqConnection)
	if err := Instance.rabbitClient.DeclareQueue(string(event.UserRegistered)); err != nil {
		log.Fatal(err)
	}
	Instance.eventBus = event.NewRabbitMqEventBus(rabbitMq.NewRabbitMqClient(rabbitMqConnection))

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

	// move to factory
	sendEmailOnUserRegistered := commands.NewSendEmailOnUserRegistered(commands.SendEmailOnUserRegisteredDependencies{
		UserRepository:      Instance.userRepository,
		EmailSender:         Instance.emailSender,
		SentEmailRepository: Instance.sentEmailRepository,
		WelcomeTemplatePath: config.Conf.Smtp.Welcome.Template,
	})
	userRegisteredEventHandler := events.NewUserRegisteredEventHandler(sendEmailOnUserRegistered, Instance.commandBus, Instance.messageRepository)
	go Instance.eventBus.Subscribe(event.UserRegistered, userRegisteredEventHandler)

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
