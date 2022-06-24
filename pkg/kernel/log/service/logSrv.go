package service

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/customError"
	"io/ioutil"
	"os"
	"time"
)

type LogSrv struct {
	filePath string
}

func NewLogSrv() *LogSrv {
	fileName := os.Getenv("LOG_FILE_PATH")

	return &LogSrv{
		filePath: fileName,
	}
}

func (l LogSrv) LogError(err error, inputTime time.Time) error {
	internalError := customError.NewInternalError(err, 2)
	err = internalError.InternalError().Error()
	file := internalError.InternalError().File()
	line := internalError.InternalError().Line()
	strTime := inputTime.Format(time.UnixDate)
	logData := fmt.Sprintf("error: %v, file: %v, line: %v, time: %v", err, file, line, strTime)

	return ioutil.WriteFile(l.filePath, []byte(logData), 0644)
}

func (l LogSrv) File() ([]byte, error) {
	return ioutil.ReadFile(l.filePath)
}
