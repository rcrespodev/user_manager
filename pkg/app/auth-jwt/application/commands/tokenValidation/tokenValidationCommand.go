package tokenValidation

type Command struct {
	token string
}

func NewCommand(token string) *Command {
	return &Command{token: token}
}
