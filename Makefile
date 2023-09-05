WEB_DIR = web
API_DIR = api
DOCKER_COMPOSE_FILE = docker-compose.yml

.PHONY: dev-web
dev-web:
	cd $(WEB_DIR) && npm run dev

.PHONY: dev-api
dev-api:
	cd $(API_DIR) && go run main.go

.PHONY: build-web
build-web:
	cd $(WEB_DIR) && npm run build

.PHONY: build-api
build-api:
	cd $(API_DIR) && go build -o bin/main

.PHONY: docker-build
docker-build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

.PHONY: docker-run
docker-run:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d 

.PHONY: docker-stop
docker-stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
