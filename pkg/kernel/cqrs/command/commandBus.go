package command

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type Bus struct {
	handlersMap HandlersMap
}

type HandlersMap map[Id]CommandHandler

func NewBus(handlersMap HandlersMap) *Bus {
	return &Bus{handlersMap: handlersMap}
}

func (b Bus) Exec(c Command, returnLog *domain.ReturnLog) {
	done := make(chan bool)
	commandId := c.BaseCommand().CommandId()
	handler, ok := b.handlersMap[commandId]
	if !ok {
		returnLog.LogError(domain.NewErrorCommand{
			Error: fmt.Errorf("there not any CommandHandler associate to commandId %v", commandId),
		})
		return
	}

	go handler.Handle(c, returnLog, done)
	if <-done == true {
		return
	}
}
