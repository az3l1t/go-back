APP_NAME = go-back
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
DOCKER_COMPOSE_FILE = docker-compose.yml

.PHONY: build run test clean

build:
	@echo "Building the application..."
	go build -o $(APP_NAME) .

test: docker-up-test wait-for-db start-app wait-for-app
	@echo "Running tests..."
	go test -v ./tests
	@$(MAKE) docker-down

start-app:
	@echo "Starting the application..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d app

wait-for-app:
	@echo "Waiting for the application to be ready..."
	@sleep 5

docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

docker-up-test:
	@echo "Starting the test database..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build test_db -d

wait-for-db:
	@echo "Waiting for the test database to be ready..."
	@while ! docker exec -it $(shell docker-compose ps -q test_db) pg_isready -U postgres; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done

docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-restart:
	docker-compose down --remove-orphans
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

