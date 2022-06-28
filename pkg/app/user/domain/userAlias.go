package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/utils/stringValidations"
	"strconv"
)

const (
	aliasMaxLenChar = 30
)

type UserAlias struct {
	alias string
	log   *returnLog.ReturnLog
}

func NewUserAlias(alias string, log *returnLog.ReturnLog) *UserAlias {
	userAlias := &UserAlias{
		alias: alias,
		log:   log,
	}

	userAlias.checkLen()
	if userAlias.log.Error() != nil {
		return nil
	}

	userAlias.checkSpecialChars()
	if userAlias.log.Error() != nil {
		return nil
	}

	return userAlias
}

func (u *UserAlias) checkLen() {
	if len(u.alias) > aliasMaxLenChar {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 006,
				Variables: message.Variables{"alias", strconv.Itoa(aliasMaxLenChar)},
			},
		})
	}
}

func (u *UserAlias) checkSpecialChars() {
	if contain, specialChars := stringValidations.ContainSpecialChars(u.alias); contain == true {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 007,
				Variables: message.Variables{"alias", specialChars[0]},
			},
		})
	}
}

func (u UserAlias) Alias() string {
	return u.alias
}
