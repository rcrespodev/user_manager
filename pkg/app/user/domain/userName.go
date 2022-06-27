package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"unicode"
)

type UserName struct {
	name string
}

func NewUserName(name string, log *returnLog.ReturnLog) *UserName {
	userName := &UserName{name: name}
	userName.parseFirstCharacterToUpperCase()
	return userName
}

func (u *UserName) parseFirstCharacterToUpperCase() {
	runeStr := []rune(u.name)
	runeStr[0] = unicode.ToUpper(runeStr[0])
	u.name = string(runeStr)
}

func (u UserName) Name() string {
	return u.name
}
