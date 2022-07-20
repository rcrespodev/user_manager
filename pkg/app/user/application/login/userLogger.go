package login

import returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type UserLogger struct {
}

func (u UserLogger) Exec(cmd *LoginUserCommand, log *returnLog.ReturnLog) {

}
