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
	Uuid       string `json:"uuid"`
	Alias      string `json:"alias"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
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
