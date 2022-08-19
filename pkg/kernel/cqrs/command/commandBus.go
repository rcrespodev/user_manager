package command

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type Bus interface {
	Exec(c Command, returnLog *domain.ReturnLog)
}
