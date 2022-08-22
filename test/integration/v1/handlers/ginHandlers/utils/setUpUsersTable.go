package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

func TableUsersSetup(commands []*domain.NewUserCommand, repository domain.UserRepository) error {
	const defaultPkf = "user"

	for _, command := range commands {
		cmdUuid, err := uuid.Parse(command.Uuid)
		if err != nil {
			return err
		}

		retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), defaultPkf)
		user := domain.NewUser(*command, retLog)
		repository.SaveUser(user, retLog)
		if retLog.Error() != nil {
			return fmt.Errorf("%v", retLog.Error())
		}
	}

	return nil
}
