test_envs:
	export APP_HOST=localhost
	export APP_PORT=8080
	export HOME_PROJECT=/home/rcrespo/github.com/rcrespodev
	export GO111MODULE=on

run_server:
	go run ./cmd/web/main.go

test_server: test_envs run_server

go_tests:
	go test ./...

test: test_envs go_tests
