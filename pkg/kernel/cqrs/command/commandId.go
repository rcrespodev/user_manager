package command

type Id uint8

const (
	RegisterUser Id = 0 + iota
	UpdateUser
	DeleteUser
)
