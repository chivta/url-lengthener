.PHONY: dev test lint build migrate tf-plan tf-apply

dev:
	docker-compose up

test:
	cd api && go test ./...
	cd frontend && npm test

lint:
	cd api && golangci-lint run
	cd frontend && npm run lint

build:
	docker-compose build

migrate:
	cd api && go run ./cmd/server/main.go

tf-plan:
	cd infra && terraform plan

tf-apply:
	cd infra && terraform apply
