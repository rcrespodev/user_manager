package userLoggedOut

type Command struct {
	userUuid string
}

func NewCommand(userUuid string) *Command {
	return &Command{userUuid: userUuid}
}

func (c Command) UserUuid() string {
	return c.userUuid
}
