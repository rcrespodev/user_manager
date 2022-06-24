package log

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/pkg/kernel/log/service"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogSrv(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	homeProject := os.Getenv("HOME_PROJECT")
	logDir := os.Getenv("LOG_FILE_PATH")

	type args struct {
		error error
		time  time.Time
	}
	type want struct {
		file string
		str  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base test",
			args: args{
				error: fmt.Errorf("test error"),
				time:  time.Date(2022, time.June, 24, 07, 22, 13, 430701331, time.Local),
			},
			want: want{
				file: "wantFile.log",
				str:  fmt.Sprintf("error: test error, file: %v/user_manager/test/unit/kernel/log/logSrv_test.go, line: 64, time: Fri Jun 24 07:22:13 -03 2022", homeProject),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.want.file); err == nil {
				err := os.Remove(tt.want.file)
				if err != nil {
					log.Fatal(err)
				}
			}

			if _, err := os.Stat(logDir); err == nil {
				err := os.Remove(logDir)
				if err != nil {
					log.Fatal(err)
				}
			}

			logSrv := service.NewLogSrv()
			if err := logSrv.LogError(tt.args.error, tt.args.time); err != nil {
				log.Fatal(err)
			}

			gotFile, err := logSrv.File()
			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile(tt.want.file, []byte(tt.want.str), 0644)
			if err != nil {
				log.Fatal(err)
			}

			wantFile, err := ioutil.ReadFile(tt.want.file)
			if err != nil {
				log.Fatal(err)
			}

			if string(gotFile) != string(wantFile) {
				t.Errorf("File()\n\t- got: %v\n\t- want: %v", string(gotFile), string(wantFile))
			}
		})
	}
}
