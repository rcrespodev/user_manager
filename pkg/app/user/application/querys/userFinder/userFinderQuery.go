package userFinder

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
)

type Query struct {
	queryArgs []domain.FindUserQuery
	baseQuery *query.BaseQuery
}

func (q Query) BaseQuery() *query.BaseQuery {
	return q.baseQuery
}

func NewQuery(queryArgs []domain.FindUserQuery) *Query {
	return &Query{
		queryArgs: queryArgs,
		baseQuery: query.NewBaseQuery(query.FindUser),
	}
}

func (q Query) QueryArgs() []domain.FindUserQuery {
	return q.queryArgs
}
