package command

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type Bus struct {
	handlersMap map[Id]Handler
}

func NewBus(handlersMap map[Id]Handler) *Bus {
	return &Bus{handlersMap: handlersMap}
}

func (b Bus) Exec(c Command, returnLog *domain.ReturnLog) {
	done := make(chan bool)
	commandId := c.CommandId()
	handler, ok := b.handlersMap[commandId]
	if !ok {
		returnLog.LogError(domain.NewErrorCommand{
			Error: fmt.Errorf("there not any Handler associate to commandId %v", commandId),
		})
		return
	}

	go handler.Handle(c, returnLog, done)
	if <-done == true {
		return
	}
}
