package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/utils/stringValidations"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"unicode"
)

type UserPassword struct {
	stringPassword string
	hashPassword   []byte
	log            *returnLog.ReturnLog
}

func NewUserPassword(password string, log *returnLog.ReturnLog) *UserPassword {
	userPassword := &UserPassword{
		stringPassword: password,
		hashPassword:   nil,
		log:            log,
	}

	userPassword.checkSpecialCharacter()
	if userPassword.log.Error() != nil {
		return nil
	}

	userPassword.unicodeCheck()
	if userPassword.log.Error() != nil {
		return nil
	}

	userPassword.checkLen()
	if userPassword.log.Error() != nil {
		return nil
	}

	userPassword.newHashPassword()
	if userPassword.log.Error() != nil {
		return nil
	}

	return userPassword
}

func (u *UserPassword) newHashPassword() {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.stringPassword), bcrypt.DefaultCost)
	if err != nil {
		u.log.LogError(returnLog.NewErrorCommand{
			Error:             err,
			NewMessageCommand: nil,
		})
		return
	}
	u.hashPassword = hashPassword
}

func (u *UserPassword) checkLen() {
	const (
		maxLen = 16
		minLen = 8
	)
	if len(u.stringPassword) > maxLen {
		//log.Printf("pass: %v", u.stringPassword)
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 6,
				Variables: message.Variables{"password", strconv.Itoa(maxLen)},
			},
		})
	}
	if len(u.stringPassword) < minLen {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 10,
				Variables: message.Variables{"password", strconv.Itoa(minLen)},
			},
		})
	}
}

func (u *UserPassword) unicodeCheck() {
	var (
		number    bool
		upperCase bool
		lowerCase bool
	)

	for _, char := range u.stringPassword {
		if ok := unicode.IsLower(char); ok {
			lowerCase = ok
		}
		if ok := unicode.IsUpper(char); ok {
			upperCase = ok
		}
		if ok := unicode.IsNumber(char); ok {
			number = ok
		}
		if number == true && upperCase == true && lowerCase == true {
			break
		}
	}

	if !number {
		u.log.LogError(returnLog.NewErrorCommand{
			Error:             nil,
			NewMessageCommand: &message.NewMessageCommand{MessageId: 9},
		})
	}

	if !upperCase {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 11,
			},
		})
	}

	if !lowerCase {
		u.log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 12,
			},
		})
	}
}

func (u *UserPassword) checkSpecialCharacter() {
	if ok, _ := stringValidations.ContainSpecialChars(u.stringPassword); !ok {
		u.log.LogError(returnLog.NewErrorCommand{
			Error:             nil,
			NewMessageCommand: &message.NewMessageCommand{MessageId: 8},
		})
	}
}

func (u UserPassword) String() string {
	return u.stringPassword
}

func (u UserPassword) Hash() []byte {
	return u.hashPassword
}
