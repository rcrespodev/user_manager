package query

import (
	"github.com/google/uuid"
	"time"
)

type QueryInterface interface {
	BaseQuery() *BaseQuery
}

type BaseQuery struct {
	queryUuid  uuid.UUID
	queryId    Id
	occurredOn time.Time
}

type Id int8

const (
	FindUser Id = 0 + iota
)

func NewBaseQuery(queryId Id) *BaseQuery {
	return &BaseQuery{
		queryId: queryId,
	}
}

func (b BaseQuery) QueryUuid() uuid.UUID {
	return b.queryUuid
}

func (b BaseQuery) QueryId() Id {
	return b.queryId
}

func (b BaseQuery) OccurredOn() time.Time {
	return b.occurredOn
}
