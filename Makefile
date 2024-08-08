SERVICE_NAME=innopolis_go_crud
PORT=18001

.PHONY: default
default: build up

.PHONY: build
build:
	docker-compose build

# Запуск контейнеров в фоне
.PHONY: up
up:
	docker-compose up -d

.PHONY: init
init: build up

.PHONY: exec
exec:
	@container_id=$$(docker ps -q --filter "name=$(SERVICE_NAME)"); \
	if [ -n "$$container_id" ]; then \
		docker exec -it $$container_id sh; \
	else \
		echo "Service $(SERVICE_NAME) is not running."; \
	fi

.PHONY: stop
stop:
	docker-compose stop

.PHONY: clean
clean:
	docker-compose down --rmi all --volumes --remove-orphans

.PHONY: check
check:
	@status_code=$$(curl -o /dev/null -s -w "%{http_code}" -X GET http://localhost:$(PORT)/ping); \
	echo "HTTP Status Code: $$status_code"; \
	if [ "$$status_code" -eq 200 ]; then \
		echo "Service is running and accessible on port $(PORT)."; \
	else \
		echo "Service is not accessible on port $(PORT)."; \
		exit 1; \
	fi
