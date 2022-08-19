package registerUser

import (
	userFinderPkg "github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"sync"
)

type Service struct {
	userRepository domain.UserRepository
	user           *domain.User
	userByAlias    *domain.User
	userByEmail    *domain.User
}

func NewService(repository domain.UserRepository) *Service {
	return &Service{userRepository: repository}
}

func (u *Service) Exec(command Command, log *returnLog.ReturnLog) {
	u.user = domain.NewUser(domain.NewUserCommand{
		Uuid:       command.uuid,
		Alias:      command.alias,
		Name:       command.name,
		SecondName: command.secondName,
		Email:      command.email,
		Password:   command.password,
	}, log)

	if u.user == nil {
		return
	}

	userFinder := userFinderPkg.NewUserFinder(u.userRepository)
	newLog := *log
	sourceUser := userFinder.Exec([]domain.FindUserQuery{
		{
			Log: &newLog,
			Where: []domain.WhereArgs{
				{
					Field: "alias",
					Value: u.user.Alias().Alias(),
				},
			},
		},
		{
			Log: &newLog,
			Where: []domain.WhereArgs{
				{
					Field: "email",
					Value: u.user.Email().Address(),
				},
			},
		},
	}, &newLog)

	if sourceUser != nil {
		var variables message.Variables
		if sourceUser.Alias == u.user.Alias().Alias() {
			variables = message.Variables{"alias", u.user.Alias().Alias()}
		}
		if sourceUser.Email == u.user.Email().Address() {
			variables = message.Variables{"email", u.user.Email().Address()}
		}
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 14,
				Variables: variables,
			},
		})
		return
	}

	u.userRepository.SaveUser(u.user, log)
	if log.Error() != nil {
		return
	}

	log.LogSuccess(&message.NewMessageCommand{
		MessageId: 1,
		Variables: message.Variables{u.user.Alias().Alias()},
	})
	return
}

func (u *Service) finUserByAlias(wg *sync.WaitGroup, log *returnLog.ReturnLog) {
	userSchema := u.userRepository.FindUser(domain.FindUserQuery{
		Log: log,
		Where: []domain.WhereArgs{
			{
				Field: "alias",
				Value: u.user.Alias().Alias(),
			},
		},
	})
	switch userSchema {
	case nil:
		u.userByAlias = nil
	default:
		u.userByAlias = domain.NewUser(domain.NewUserCommand{
			Uuid:       userSchema.Uuid,
			Alias:      userSchema.Alias,
			Name:       userSchema.Name,
			SecondName: userSchema.SecondName,
			Email:      userSchema.Email,
			Password:   u.user.Password().String(),
			IgnorePass: false,
		}, log)
	}
	wg.Done()
}

func (u *Service) finUserByEmail(wg *sync.WaitGroup, log *returnLog.ReturnLog) {
	userSchema := u.userRepository.FindUser(domain.FindUserQuery{
		Log: log,
		Where: []domain.WhereArgs{
			{
				Field: "email",
				Value: u.user.Email().Address(),
			},
		},
	})
	switch userSchema {
	case nil:
		u.userByEmail = nil
	default:
		u.userByEmail = domain.NewUser(domain.NewUserCommand{
			Uuid:       userSchema.Uuid,
			Alias:      userSchema.Alias,
			Name:       userSchema.Name,
			SecondName: userSchema.SecondName,
			Email:      userSchema.Email,
			Password:   u.user.Password().String(),
			IgnorePass: false,
		}, log)
	}
	wg.Done()
}
