package userFinder

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"sync"
)

type UserFinder struct {
	UserRepository domain.UserRepository
}

func NewUserFinder(userRepository domain.UserRepository) *UserFinder {
	return &UserFinder{UserRepository: userRepository}
}

func (u UserFinder) Exec(queryArgs []domain.FindUserQuery, log *returnLog.ReturnLog) *domain.UserSchema {
	userChan := make(chan *domain.UserSchema)
	done := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	for _, query := range queryArgs {
		go u.findUser(query, userChan)
	}

	var userSchema *domain.UserSchema
	for i := range queryArgs {
		user := <-userChan
		if user != nil && userSchema == nil {
			userSchema = user
			go func() {
				wg.Done()
			}()
			break
		}

		//if last iteration and user not found yet, set wg = done
		if len(queryArgs) == (i + 1) {
			go func() {
				wg.Done()
			}()
		}
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	<-done

	if userSchema == nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				ObjectId:   "",
				MessageId:  17,
				MessagePkg: "user",
				Variables:  message.Variables{},
			},
		})
	}

	return userSchema
}

func (u UserFinder) findUser(query domain.FindUserQuery, userChan chan *domain.UserSchema) {
	userChan <- u.UserRepository.FindUser(query)
}
