package command

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type BusInterface interface {
	Exec(c Command, returnLog *domain.ReturnLog)
}
