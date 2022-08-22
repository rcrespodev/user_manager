package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"time"
)

type SentEmailRepository interface {
	Save(schema *SentEmailSchema, log *returnLog.ReturnLog)
	Get(userUuid string) []*SentEmailSchema
}

type SentEmailSchema struct {
	UserUuid string
	Sent     bool
	SentOn   time.Time
	Error    string
}
