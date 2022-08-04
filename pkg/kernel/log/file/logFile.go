package file

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogFile type receives internal errors and write YYYY-MM-DD-errors.log file
// with the errors received.
// The method LogInternalError() is called each time ReturnLog type execute method LogError() and command.Error != nil
// The log file path are defined in LOG_FILE_PATH env.
type LogFile struct {
	logFilePath string
}

type LogInternalErrorCommand struct {
	Error error
	File  string
	Line  string
	Time  time.Time
}

func NewLogFile(logFilePath string) *LogFile {
	return &LogFile{
		logFilePath: logFilePath,
	}
}

func (l *LogFile) LogInternalError(command LogInternalErrorCommand) {
	l.writeFile(command)
}

func (l *LogFile) writeFile(command LogInternalErrorCommand) {
	fileName := fmt.Sprintf("%v-%v-%v-errors.log",
		command.Time.Year(), command.Time.Month(), command.Time.Day())

	path := fmt.Sprintf("%s/%s", l.logFilePath, fileName)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	t := fmt.Sprintf("%v:%v:%v", command.Time.Hour(), command.Time.Minute(), command.Time.Second())
	appendLine := fmt.Sprintf("Time: %s; Error: %s; File: %s; Line: %v\n",
		t, command.Error.Error(), command.File, command.Line)
	if _, err := f.Write([]byte(appendLine)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
