package register

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"sync"
)

type UserRegistration struct {
	userRepository domain.UserRepository
}

func NewUserRegistration(repository domain.UserRepository) *UserRegistration {
	return &UserRegistration{userRepository: repository}
}

func (u *UserRegistration) Exec(command RegisterUserCommand, log *returnLog.ReturnLog) {
	user := domain.NewUser(domain.NewUserCommand{
		Uuid:       command.uuid,
		Alias:      command.alias,
		Name:       command.name,
		SecondName: command.secondName,
		Email:      command.email,
		Password:   command.password,
	}, log)

	if user == nil {
		return
	}

	waitGroup := &sync.WaitGroup{}
	var userAlias *domain.User
	var userEmail *domain.User
	waitGroup.Add(2)
	go func() {
		userAlias = u.userRepository.FindUserByAlias(domain.FindByAliasCommand{
			Alias: user.Alias(),
			FindUserCommand: domain.FindUserCommand{
				Password: user.Password().String(),
				Log:      log,
				Wg:       waitGroup,
			},
		})
		//userAlias = u.userRepository.FindUserByAlias(user.Alias(), log, waitGroup)
	}()
	go func() {
		userEmail = u.userRepository.FindUserByEmail(domain.FindByEmailCommand{
			Email: user.Email(),
			FindUserCommand: domain.FindUserCommand{
				Password: user.Password().String(),
				Log:      log,
				Wg:       waitGroup,
			},
		})
		//userEmail = u.userRepository.FindUserByEmail(user.Email(), log, waitGroup)
	}()
	waitGroup.Wait()

	if log.Error() != nil {
		return
	}

	if userAlias != nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 14,
				Variables: message.Variables{"alias", user.Alias().Alias()},
			},
		})
		return
	}

	if userEmail != nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 14,
				Variables: message.Variables{"email", user.Email().Address()},
			},
		})
		return
	}

	waitGroup.Add(1)
	u.userRepository.SaveUser(user, log, waitGroup)
	waitGroup.Wait()

	if log.Error() != nil {
		return
	}

	log.LogSuccess(&message.NewMessageCommand{
		MessageId: 1,
		Variables: message.Variables{user.Alias().Alias()},
	})
	return
}
