package command

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type Handler interface {
	Handle(command Command, returnLog *domain.ReturnLog, done chan bool)
}
