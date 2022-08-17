package commands

import (
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"time"
)

type SendEmailOnUserRegistered struct {
	userRepository      domain.UserRepository
	emailSender         emailSenderDomain.EmailSender
	sentEmailRepository emailSenderDomain.SentEmailRepository
	welcomeTemplatePath string
}

type SendEmailOnUserRegisteredDependencies struct {
	UserRepository      domain.UserRepository
	EmailSender         emailSenderDomain.EmailSender
	SentEmailRepository emailSenderDomain.SentEmailRepository
	WelcomeTemplatePath string
}

func NewSendEmailOnUserRegistered(cmd SendEmailOnUserRegisteredDependencies) *SendEmailOnUserRegistered {
	return &SendEmailOnUserRegistered{
		userRepository:      cmd.UserRepository,
		emailSender:         cmd.EmailSender,
		welcomeTemplatePath: cmd.WelcomeTemplatePath,
		sentEmailRepository: cmd.SentEmailRepository,
	}
}

func (s SendEmailOnUserRegistered) Exec(command *SendEmailOnUserRegisteredCommand, log *returnLog.ReturnLog) {
	userUuid := command.BaseCommand().AggregateId().String()
	log.SetObjectId(userUuid)

	userSchema := s.userRepository.FindUser(domain.FindUserQuery{
		Log: log,
		Where: []domain.WhereArgs{
			{
				Field: "uuid",
				Value: userUuid,
			},
		},
	})

	if userSchema == nil || log.Error() != nil {
		return
	}

	user := domain.NewUser(domain.NewUserCommand{
		Uuid:       userSchema.Uuid,
		Alias:      userSchema.Alias,
		Name:       userSchema.Name,
		SecondName: userSchema.SecondName,
		Email:      userSchema.Email,
		Password:   "",
		IgnorePass: true,
	}, log)
	if log.Error() != nil {
		return
	}

	welcomeEmail := emailSenderDomain.NewWelcomeEmail(user, s.welcomeTemplatePath, log)
	if log.Error() != nil {
		return
	}

	sentEmailSchema := &emailSenderDomain.SentEmailSchema{
		UserUuid: user.Uuid().String(),
		SentOn:   time.Now(),
	}

	s.emailSender.SendEmail(&emailSenderDomain.SendEmailCommand{
		To:      welcomeEmail.To(),
		Subject: welcomeEmail.Subject(),
		Body:    welcomeEmail.Body(),
	}, log)
	if log.Error() != nil {
		sentEmailSchema.Error = log.Error().InternalError().Error().Error()
		s.sentEmailRepository.Save(sentEmailSchema, log)
		return
	}

	sentEmailSchema.Sent = true
	s.sentEmailRepository.Save(sentEmailSchema, log)
	if log.Error() != nil {
		return
	}

	log.LogSuccess(&message.NewMessageCommand{
		MessageId:  0,
		MessagePkg: "email_sender",
		Variables:  message.Variables{user.Email().Address()},
	})
	return
}
