package loginUser

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type Command struct {
	aliasOrEmail string
	password     string
	baseCommand  *command.BaseCommand
}

type ClientArgs struct {
	AliasOrEmail string `json:"alias_or_email"`
	Password     string `json:"password"`
}

func NewLoginUserCommand(args ClientArgs) *Command {
	return &Command{
		aliasOrEmail: args.AliasOrEmail,
		password:     args.Password,
		baseCommand:  command.NewBaseCommand(args.AliasOrEmail, command.LoginUser),
	}
}

func (l Command) Password() string {
	return l.password
}

func (l Command) AliasOrEmail() string {
	return l.aliasOrEmail
}

func (l Command) BaseCommand() *command.BaseCommand {
	return l.baseCommand
}
