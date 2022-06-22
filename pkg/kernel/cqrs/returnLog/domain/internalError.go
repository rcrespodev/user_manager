package domain

type InternalError struct {
	error error
	file  string
	line  int
}

func (i InternalError) Error() error {
	return i.error
}

func (i InternalError) File() string {
	return i.file
}
