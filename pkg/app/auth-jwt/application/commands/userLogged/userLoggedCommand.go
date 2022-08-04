package userLogged

import "github.com/google/uuid"

type Command struct {
	userUuid uuid.UUID
}

func NewCommand(userUuid uuid.UUID) *Command {
	return &Command{userUuid: userUuid}
}

func (c Command) UserUuid() uuid.UUID {
	return c.userUuid
}
