package userFinder

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type QueryHandler struct {
	userFinder *UserFinder
	query      *Query
}

func NewQueryHandler(userFinder *UserFinder) *QueryHandler {
	return &QueryHandler{userFinder: userFinder}
}

func (q QueryHandler) Query(query *query.Query, log *returnLog.ReturnLog, data chan interface{}) {
	findUserQuery, ok := query.Args().(*Query)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		data <- nil
		return
	}
	q.query = findUserQuery

	data <- q.userFinder.Exec(findUserQuery.QueryArgs(), log)
}
