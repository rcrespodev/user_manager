package delete

type DeleteUserCommand struct {
	userUuid string
}

func NewDeleteUserCommand(uuid string) *DeleteUserCommand {
	return &DeleteUserCommand{userUuid: uuid}
}

func (u DeleteUserCommand) UserUuid() string {
	return u.userUuid
}
