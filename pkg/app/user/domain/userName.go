package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/utils/stringValidations"
	"strconv"
	"unicode"
)

type UserName struct {
	name string
	log  *returnLog.ReturnLog
}

func NewUserName(name string, log *returnLog.ReturnLog) *UserName {
	userName := &UserName{
		name: name,
		log:  log,
	}
	userName.parseFirstCharacterToUpperCase()

	userName.checkLen()
	if userName.log.Error() != nil {
		return nil
	}

	userName.checkSpecialChars()
	if userName.log.Error() != nil {
		return nil
	}
	return userName
}

func (u *UserName) checkLen() {
	const (
		maxLen = 50
	)
	if len(u.name) > maxLen {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 6,
				Variables: message.Variables{"name", strconv.Itoa(maxLen)},
			},
		})
	}
}

func (u *UserName) parseFirstCharacterToUpperCase() {
	runeStr := []rune(u.name)
	runeStr[0] = unicode.ToUpper(runeStr[0])
	u.name = string(runeStr)
}

func (u *UserName) checkSpecialChars() {
	if contain, specialChars := stringValidations.ContainSpecialChars(u.name); contain {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 007,
				Variables: message.Variables{"name", specialChars[0]},
			},
		})
	}
}

func (u UserName) Name() string {
	return u.name
}
