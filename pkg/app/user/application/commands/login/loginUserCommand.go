package login

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type LoginUserCommand struct {
	aliasOrEmail string
	password     string
	baseCommand  *command.BaseCommand
}

type ClientArgs struct {
	AliasOrEmail string `json:"alias_or_email"`
	Password     string `json:"password"`
}

func NewLoginUserCommand(args ClientArgs) *LoginUserCommand {
	return &LoginUserCommand{
		aliasOrEmail: args.AliasOrEmail,
		password:     args.Password,
		baseCommand:  command.NewBaseCommand(args.AliasOrEmail, command.LoginUser),
	}
}

func (l LoginUserCommand) Password() string {
	return l.password
}

func (l LoginUserCommand) AliasOrEmail() string {
	return l.aliasOrEmail
}

func (l LoginUserCommand) BaseCommand() *command.BaseCommand {
	return l.baseCommand
}
