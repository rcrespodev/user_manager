package customError

import (
	"runtime"
	"strings"
)

type InternalError struct {
	error error
	file  string
	line  int
}

func (i *InternalError) parsePath() {
	sep := "/"
	if runtime.GOOS == "windows" {
		sep = "\\"
	}
	splitPath := strings.Split(i.file, sep)
	for index, s := range splitPath {
		if s == "user_manager" { // root project
			i.file = strings.Join(splitPath[index:], sep)
			break
		}
	}
}

func (i *InternalError) Error() error {
	return i.error
}

func (i *InternalError) File() string {
	return i.file
}

func (i *InternalError) Line() int {
	return i.line
}
