.PHONY:test run lint fmt

IMAGE_NAME=web-page-analyzer
VERSION=latest

run:
	go run ./main.go

lint:
	@$(shell go env GOPATH)/bin/staticcheck ./...

fmt:
	go fmt

test:
	go test -v ./... -cover 

build:
	docker build -t $(IMAGE_NAME):$(VERSION) .

run-docker:
	docker run -d -p 80:80 --name web-page-analyzer $(IMAGE_NAME):$(VERSION)