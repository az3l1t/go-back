APP_NAME = go-back
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
DOCKER_COMPOSE_FILE = docker-compose.yml

.PHONY: build run test clean docker-up docker-down wait-for-db wait-for-app start-app

build:
	@echo "Building the application..."
	go build -o $(APP_NAME) .

test: docker-up wait-for-db start-app wait-for-app
	@echo "Running tests..."
	# Запуск тестов в фоновом режиме, подключаясь к запущенному приложению и базе данных
	go test -v ./tests
	@$(MAKE) docker-down

start-app:
	@echo "Starting the application in the background..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d app

wait-for-app:
	@echo "Waiting for the application to be ready..."
	@sleep 10

docker-up:
	@echo "Starting database and application containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d db

wait-for-db:
	@echo "Waiting for the database to be ready..."
	@while ! docker exec -it $(shell docker-compose ps -q db) pg_isready -U postgres; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done

docker-down:
	@echo "Stopping and removing all containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --remove-orphans

docker:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(APP_NAME)
