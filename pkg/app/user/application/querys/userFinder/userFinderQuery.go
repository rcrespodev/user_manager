package userFinder

import "github.com/rcrespodev/user_manager/pkg/app/user/domain"

type Query struct {
	queryArgs []domain.FindUserQuery
}

func NewQuery(queryArgs []domain.FindUserQuery) *Query {
	return &Query{queryArgs: queryArgs}
}

func (q Query) QueryArgs() []domain.FindUserQuery {
	return q.queryArgs
}
