package delete

type DeleteUserCommand struct {
	userUuid string
}

type ClientArgs struct {
	UserUuid string
}

func NewDeleteUserCommand(uuid string) *DeleteUserCommand {
	return &DeleteUserCommand{userUuid: uuid}
}

func (u DeleteUserCommand) UserUuid() string {
	return u.userUuid
}
