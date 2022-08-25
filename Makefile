run: stop_app up_app_d

run_online: stop_app up_app_online

stop_app:
	docker-compose -f docker-compose.yaml stop

up_app_d:
	docker-compose -f docker-compose.yaml up -d --build

up_app_online:
	docker-compose -f docker-compose.yaml up --build

stop_services:
	sudo service redis stop | sudo service mysql stop

stop_test_services:
	docker stop test_app_mysql test_app_redis test_app_rabbitmq

go_tests:
	export GO111MODULE=on && go test -v ./...

cert:
	openssl genrsa -out cert/id_rsa 4096
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub

run_demo:
	export GO111MODULE=on && go run ./demo/http/cmd/main.go
