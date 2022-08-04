package query

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type QueryHandler interface {
	Query(query *Query, log *domain.ReturnLog, data chan interface{})
}
