package file

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/log/file"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogFile(t *testing.T) {
	type args struct {
		commands []file.LogInternalErrorCommand
		filePath string
		fileName string
	}
	type want struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base test",
			args: args{
				commands: []file.LogInternalErrorCommand{
					{
						Error: fmt.Errorf("test error 1"),
						File:  "app/pkg/kernel/kernel.go",
						Line:  "50",
						Time:  time.Date(2022, time.June, 20, 9, 55, 43, 0, time.Local),
					},
					{
						Error: fmt.Errorf("test error 2"),
						File:  "app/pkg/kernel/utils/stringValidations/specialCharacters.go",
						Line:  "4",
						Time:  time.Date(2022, time.June, 20, 18, 30, 5, 0, time.Local),
					},
					{
						Error: fmt.Errorf("test error 3"),
						File:  "app/pkg/server/serve.go",
						Line:  "18",
						Time:  time.Date(2022, time.June, 20, 23, 47, 0, 0, time.Local),
					},
				},
				filePath: "gotFiles",
				fileName: "2022-June-20-errors.log",
			},
			want: want{
				filePath: "wantFiles/2022-June-20-errors.log",
			},
		},
	}
	for _, tt := range tests {
		gotFilePath := fmt.Sprintf("%s/%s", tt.args.filePath, tt.args.fileName)

		// delete file if exists
		if _, err := os.Stat(gotFilePath); err == nil {
			err := os.Remove(gotFilePath)
			if err != nil {
				log.Fatal(err)
			}
		}

		logFileSrv := file.NewLogFile(tt.args.filePath)
		for _, command := range tt.args.commands {
			logFileSrv.LogInternalError(file.LogInternalErrorCommand{
				Error: command.Error,
				File:  command.File,
				Line:  command.Line,
				Time:  command.Time,
			})
		}
		wantFile, err := os.ReadFile(tt.want.filePath)
		require.NoError(t, err)

		gotFile, err := os.ReadFile(gotFilePath)
		require.NoError(t, err)

		require.EqualValues(t, wantFile, gotFile)
	}
}
