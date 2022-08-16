package register

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type RegisterUserCommand struct {
	uuid        string
	alias       string
	name        string
	secondName  string
	email       string
	password    string
	baseCommand *command.BaseCommand
}

type ClientArgs struct {
	Uuid       string `json:"uuid"`
	Alias      string `json:"alias"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func NewRegisterUserCommand(args ClientArgs) *RegisterUserCommand {
	return &RegisterUserCommand{
		uuid:        args.Uuid,
		alias:       args.Alias,
		name:        args.Name,
		secondName:  args.SecondName,
		email:       args.Email,
		password:    args.Password,
		baseCommand: command.NewBaseCommand(args.Uuid, command.RegisterUser),
	}
}

func (r RegisterUserCommand) BaseCommand() *command.BaseCommand {
	return r.baseCommand
}
