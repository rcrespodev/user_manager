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

run_test:
	docker-compose -f test/docker-compose.yaml up --build --abort-on-container-exit
	#docker-compose -f test/docker-compose.yml down --volumes
	docker-compose -f test/docker-compose.yaml down --volumes