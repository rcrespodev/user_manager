package query

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type Bus struct {
	handlersMap HandlersMap
}

type HandlersMap map[Id]QueryHandler

func NewBus(handlersMap HandlersMap) *Bus {
	return &Bus{handlersMap: handlersMap}
}

func (b Bus) Exec(q *Query, returnLog *domain.ReturnLog) (data interface{}) {
	dataCh := make(chan interface{})
	queryId := q.QueryId()
	handler, ok := b.handlersMap[queryId]
	if !ok {
		returnLog.LogError(domain.NewErrorCommand{
			Error: fmt.Errorf("there not any CommandHandler associate to queryId %v", queryId),
		})
		return
	}

	go handler.Query(q, returnLog, dataCh)
	return <-dataCh
}
