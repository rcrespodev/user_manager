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

	//var userAlias chan *domain.User
	//var userEmail chan *domain.User
	//go u.userRepository.FindUserByAlias(domain.FindByAliasCommand{
	//	Alias: user.Alias(),
	//	FindUserCommand: domain.FindUserCommand{
	//		Password: user.Password().String(),
	//		Log:      log,
	//		//Wg:       wg,
	//	},
	//}, userChan)
	//wg.Add(2)
	//go func() {
	//u.userRepository.FindUserByAlias(domain.FindByAliasCommand{
	//	Alias: user.Alias(),
	//	FindUserCommand: domain.FindUserCommand{
	//		Password: user.Password().String(),
	//		Log:      log,
	//		//Wg:       wg,
	//	},
	//}, userAlias)
	////userAlias = u.userRepository.FindUserByAlias(user.Alias(), log, wg)
	//}()
	//go func() {
	//go u.userRepository.FindUserByEmail(domain.FindByEmailCommand{
	//	Email: user.Email(),
	//	FindUserCommand: domain.FindUserCommand{
	//		Password: user.Password().String(),
	//		Log:      log,
	//		//Wg:       wg,
	//	},
	//}, userChan)
	//userEmail = u.userRepository.FindUserByEmail(user.Email(), log, wg)
	//}()
	//wg.Wait()

	//if log.Error() != nil {
	//	return
	//}

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

	//go u.userRepository.FindUser(domain.FindUserCommand{
	//	Password: u.user.Password().String(),
	//	Log:      log,
	//	Wg:       wg,
	//	Where: []domain.WhereArgs{
	//		{
	//			Field: "email",
	//			Value: u.user.Email().Address(),
	//		},
	//	},
	//})
	//
	//wg.Wait()

	//users := make([]*domain.User, 2)
	//for i, _ := range users {
	//	users[i] = <-userChan
	//	if users[i] != nil {
	//		log.LogError(returnLog.NewErrorCommand{
	//			NewMessageCommand: u.userExistsMessage(user, users[i]),
	//		})
	//		//break
	//	}
	//}
	//if log.Error() != nil {
	//	return
	//}

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
	u.userByAlias = u.userRepository.FindUser(domain.FindUserCommand{
		Password: u.user.Password().String(),
		Log:      log,
		//Wg:       wg,
		Where: []domain.WhereArgs{
			{
				Field: "alias",
				Value: u.user.Alias().Alias(),
			},
		},
	})
	wg.Done()
}

func (u *UserRegistration) finUserByEmail(wg *sync.WaitGroup, log *returnLog.ReturnLog) {
	u.userByEmail = u.userRepository.FindUser(domain.FindUserCommand{
		Password: u.user.Password().String(),
		Log:      log,
		//Wg:       waitGroup,
		Where: []domain.WhereArgs{
			{
				Field: "email",
				Value: u.user.Email().Address(),
			},
		},
	})
	wg.Done()
}

func (u UserRegistration) userExistsMessage(userA, userB *domain.User) *message.NewMessageCommand {
	if userA.Email() == userB.Email() {
		return &message.NewMessageCommand{
			MessageId: 14,
			Variables: message.Variables{"email", userA.Email().Address()},
		}
	}
	if userA.Alias() == userB.Alias() {
		return &message.NewMessageCommand{
			MessageId: 14,
			Variables: message.Variables{"alias", userA.Alias().Alias()},
		}
	}
	return nil
}
