package delete

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type UserDeleter struct {
	userRepository domain.UserRepository
}

func NewUserDeleter(userRepository domain.UserRepository) *UserDeleter {
	return &UserDeleter{userRepository: userRepository}
}

func (u *UserDeleter) Exec(command *DeleteUserCommand, log *returnLog.ReturnLog) {
	log.SetObjectId(command.UserUuid())

	userSchema := u.userRepository.FindUser(domain.FindUserQuery{
		Log: log,
		Where: []domain.WhereArgs{
			{
				Field: "uuid",
				Value: command.UserUuid(),
			},
		},
	})
	if log.Error() != nil {
		return
	}

	user := domain.NewUser(domain.NewUserCommand{
		Uuid:       userSchema.Uuid,
		Alias:      userSchema.Alias,
		Name:       userSchema.Name,
		SecondName: userSchema.SecondName,
		Email:      userSchema.Email,
		Password:   "ignore",
		IgnorePass: true,
	}, log)

	if log.Error() != nil {
		return
	}

	u.userRepository.DeleteUser(user, log)

	if log.Error() == nil {
		log.LogSuccess(&message.NewMessageCommand{
			MessageId: 3,
			Variables: message.Variables{user.Alias().Alias()},
		})
	}
}
