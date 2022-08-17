package event

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type Handler interface {
	Handle(events []Event, log *domain.ReturnLog, done chan bool)
}
