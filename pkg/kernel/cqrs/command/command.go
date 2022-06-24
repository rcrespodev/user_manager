package command

import "github.com/google/uuid"

type Command struct {
	commandId Id // examples of commandId: register_new_user, update_user
	uuid      uuid.UUID
	args      interface{}
}

func NewCommand(commandId Id, args interface{}) *Command {
	return &Command{
		commandId: commandId,
		uuid:      uuid.New(),
		args:      args,
	}
}

func (c Command) CommandId() Id {
	return c.commandId
}

func (c Command) Uuid() uuid.UUID {
	return c.uuid
}

func (c Command) Args() interface{} {
	return c.args
}
