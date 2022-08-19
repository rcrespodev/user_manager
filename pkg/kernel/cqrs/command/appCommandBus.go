package command

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type AppBus struct {
	handlersMap HandlersMap
}

type HandlersMap map[Id]CommandHandler

func NewAppBus(handlersMap HandlersMap) *AppBus {
	return &AppBus{handlersMap: handlersMap}
}

func (b AppBus) Exec(c Command, returnLog *domain.ReturnLog) {
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
