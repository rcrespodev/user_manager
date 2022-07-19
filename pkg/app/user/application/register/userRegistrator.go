package register

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"sync"
)

type UserRegistration struct {
	userRepository domain.UserRepository
	user           *domain.User
	userByAlias    *domain.User
	userByEmail    *domain.User
}

func NewUserRegistration(repository domain.UserRepository) *UserRegistration {
	return &UserRegistration{userRepository: repository}
}

func (u *UserRegistration) Exec(command RegisterUserCommand, log *returnLog.ReturnLog) {
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

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	wg.Add(1)
	go u.finUserByAlias(wg, log)
	wg.Add(1)
	go u.finUserByEmail(wg, log)

	go func() {
		wg.Wait()
		done <- true
	}()

	<-done

	if u.userByAlias != nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 14,
				Variables: message.Variables{"alias", u.user.Alias().Alias()},
			},
		})
		return
	}

	if u.userByEmail != nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 14,
				Variables: message.Variables{"email", u.user.Email().Address()},
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

func (u *UserRegistration) finUserByAlias(wg *sync.WaitGroup, log *returnLog.ReturnLog) {
	userSchema := u.userRepository.FindUser(domain.FindUserCommand{
		Password: u.user.Password().String(),
		Log:      log,
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

func (u *UserRegistration) finUserByEmail(wg *sync.WaitGroup, log *returnLog.ReturnLog) {
	userSchema := u.userRepository.FindUser(domain.FindUserCommand{
		Password: u.user.Password().String(),
		Log:      log,
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
