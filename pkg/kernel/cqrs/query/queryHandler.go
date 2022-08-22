package query

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type QueryHandler interface {
	Query(query QueryInterface, log *domain.ReturnLog, data chan interface{})
}
