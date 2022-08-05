package login

import (
	userFinderPkg "github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type UserLogger struct {
	userRepository domain.UserRepository
	userSchema     *domain.UserSchema
	log            *returnLog.ReturnLog
	password       *domain.UserPassword
	user           *domain.User
}

func NewUserLogger(repository domain.UserRepository) *UserLogger {
	return &UserLogger{
		userRepository: repository,
	}
}

func (u *UserLogger) Exec(cmd *LoginUserCommand, log *returnLog.ReturnLog) {
	var (
		login = false
	)

	u.user = nil

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

	userFinder := userFinderPkg.NewUserFinder(u.userRepository)
	u.userSchema = userFinder.Exec([]domain.FindUserQuery{
		{
			Log: u.log,
			Where: []domain.WhereArgs{
				{
					Field: "alias",
					Value: u.getAliasQueryValue(alias),
				},
			},
		},
		{
			Log: u.log,
			Where: []domain.WhereArgs{
				{
					Field: "email",
					Value: u.getEmailQueryValue(email),
				},
			},
		},
	}, u.log)

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

func (u *UserLogger) getEmailQueryValue(email *domain.UserEmail) string {
	if email != nil {
		return email.Address()
	}
	return ""
}

func (u *UserLogger) getAliasQueryValue(alias *domain.UserAlias) string {
	if alias != nil {
		return alias.Alias()
	}
	return ""
}
