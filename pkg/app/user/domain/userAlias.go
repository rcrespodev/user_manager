package domain

import returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type UserAlias struct {
	alias string
}

func NewUserAlias(alias string, log *returnLog.ReturnLog) *UserAlias {
	return &UserAlias{alias: alias}
}

func (u UserAlias) Alias() string {
	return u.alias
}
