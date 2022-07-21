package login

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"sync"
)

type UserLogger struct {
	repository domain.UserRepository
	userSchema *domain.UserSchema
	log        *returnLog.ReturnLog
	password   *domain.UserPassword
	user       *domain.User
}

func NewUserLogger(repository domain.UserRepository) *UserLogger {
	return &UserLogger{repository: repository}
}

func (u *UserLogger) Exec(cmd *LoginUserCommand, log *returnLog.ReturnLog) {
	//var login *bool
	var (
		login = false
	)

	u.log = log
	u.log.SetObjectId(cmd.aliasOrEmail)

	defer func(log *returnLog.ReturnLog, login *bool) {
		if !*login {
			log.LogError(returnLog.NewErrorCommand{
				NewMessageCommand: &message.NewMessageCommand{
					MessageId: 15,
					Variables: message.Variables{},
				},
				Overwritten: true,
			})
		}
	}(u.log, &login)

	u.password = domain.NewUserPassword(cmd.Password(), u.log)
	if u.log.Error() != nil {
		return
	}

	aliasLog := *u.log
	emailLOg := *u.log
	alias := domain.NewUserAlias(cmd.aliasOrEmail, &aliasLog)
	email := domain.NewUserEmail(cmd.aliasOrEmail, &emailLOg)
	if email == nil && alias == nil {
		return
	}

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	wg.Add(1)
	go u.loginWithAlias(alias, wg)
	wg.Add(1)
	go u.loginWithEmail(email, wg)

	go func() {
		wg.Wait()
		done <- true
	}()

	<-done

	if u.userSchema != nil {
		u.user = domain.LoginUser(u.password, u.userSchema, u.log)
	}

	if u.user != nil {
		log.LogSuccess(&message.NewMessageCommand{
			MessageId: 0,
			Variables: message.Variables{u.user.Alias().Alias()},
		})
		login = true
	}
}

func (u *UserLogger) loginWithAlias(alias *domain.UserAlias, wg *sync.WaitGroup) {
	if alias == nil {
		wg.Done()
		return
	}
	userSchema := u.repository.FindUser(domain.FindUserQuery{
		Log: u.log,
		Where: []domain.WhereArgs{
			{
				Field: "alias",
				Value: alias.Alias(),
			},
		},
	})
	if u.userSchema == nil {
		u.userSchema = userSchema
		wg.Done()
		return
	}
	wg.Done()
}

func (u *UserLogger) loginWithEmail(email *domain.UserEmail, wg *sync.WaitGroup) {
	if email == nil {
		wg.Done()
		return
	}
	userSchema := u.repository.FindUser(domain.FindUserQuery{
		Log: u.log,
		Where: []domain.WhereArgs{
			{
				Field: "email",
				Value: email.Address(),
			},
		},
	})
	if u.userSchema == nil {
		u.userSchema = userSchema
		wg.Done()
		return
	}
	wg.Done()
}
