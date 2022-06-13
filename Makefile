test_envs:
	export APP_HOST=localhost
	export APP_PORT=8080

run_server:
	go run ./cmd/web/main.go

test_server: test_envs run_server

run_tests:
	go test ./...

testing: test_envs run_tests
