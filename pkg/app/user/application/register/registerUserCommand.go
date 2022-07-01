package register

type RegisterUserCommand struct {
	uuid       string
	alias      string
	name       string
	secondName string
	email      string
	password   string
}

type ClientArgs struct {
	Uuid       string
	Alias      string
	Name       string
	SecondName string
	Email      string
	Password   string
}

func NewRegisterUserCommand(args ClientArgs) *RegisterUserCommand {
	return &RegisterUserCommand{
		uuid:       args.Uuid,
		alias:      args.Alias,
		name:       args.Name,
		secondName: args.SecondName,
		email:      args.Email,
		password:   args.Password,
	}
}
