PROJECTNAME := $(shell basename "$(PWD)")

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: storage-set
storage-set:
	@docker-compose -f ./deployment/dynamodb/compose.yaml --project-directory . up -d
	@sleep 3
	@AWS_PAGER="" aws dynamodb create-table --cli-input-json file://deployment/dynamodb/create-table.json --endpoint-url http://localhost:8000

.PHONY: storage-clean
storage-clean:
	@docker-compose -f ./deployment/dynamodb/compose.yaml --project-directory . down
	@rm -rf ./docker/dynamodb

.PHONY: service-build
service-build:
	@docker build --tag ${SERVICE_NAME}:$(shell git rev-parse HEAD) -f ./build/Dockerfile .

.PHONY: service-up
service-up:
	@docker-compose -f ./deployment/compose-local.yaml --project-directory . up -d