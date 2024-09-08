SERVICE_NAME=innopolis_go_crud
PORT=18001
JAEGER_UI_PORT=16686
JAEGER_COLLECTOR_PORT=14268
PROMETHEUS_PORT=9090

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

.PHONY: check_jaeger
check_jaeger:
	@status_code=$$(curl -o /dev/null -s -w "%{http_code}" -X GET http://localhost:$(JAEGER_UI_PORT)); \
	echo "HTTP Status Code for Jaeger UI: $$status_code"; \
	if [ "$$status_code" -eq 200 ]; then \
		echo "Jaeger UI is accessible on port $(JAEGER_UI_PORT)."; \
	else \
		echo "Jaeger UI is not accessible on port $(JAEGER_UI_PORT)."; \
		exit 1; \
	fi

.PHONY: check_prometheus
check_prometheus:
	@status_code=$$(curl -o /dev/null -s -w "%{http_code}" -X GET http://localhost:$(PROMETHEUS_PORT)/metrics); \
	echo "HTTP Status Code for Prometheus: $$status_code"; \
	if [ "$$status_code" -eq 200 ]; then \
		echo "Prometheus is running and accessible on port $(PROMETHEUS_PORT)."; \
	else \
		echo "Prometheus is not accessible on port $(PROMETHEUS_PORT)."; \
		exit 1; \
	fi
