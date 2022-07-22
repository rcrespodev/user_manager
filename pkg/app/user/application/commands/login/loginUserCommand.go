package login

type LoginUserCommand struct {
	aliasOrEmail string
	password     string
}

type ClientArgs struct {
	AliasOrEmail string `json:"alias_or_email"`
	Password     string `json:"password"`
}

func NewLoginUserCommand(args ClientArgs) *LoginUserCommand {
	return &LoginUserCommand{
		aliasOrEmail: args.AliasOrEmail,
		password:     args.Password,
	}
}

func (l LoginUserCommand) Password() string {
	return l.password
}

func (l LoginUserCommand) AliasOrEmail() string {
	return l.aliasOrEmail
}
