package deleteUser

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type Service struct {
	userRepository domain.UserRepository
}

func NewService(userRepository domain.UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (u *Service) Exec(command *Command, log *returnLog.ReturnLog) {
	log.SetObjectId(command.UserUuid())

	finder := userFinder.NewUserFinder(u.userRepository)
	userSchema := finder.Exec([]domain.FindUserQuery{
		{
			Log: log,
			Where: []domain.WhereArgs{
				{
					Field: "uuid",
					Value: command.UserUuid(),
				},
			},
		},
	}, log)

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
