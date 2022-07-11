package command

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type CommandHandler interface {
	Handle(command Command, log *domain.ReturnLog, done chan bool)
}
