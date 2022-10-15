# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=gofunction
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
DOCKER_REPO=fdocampo
DOCKER_CONTAINER=frutal

all: build-linux zip-project

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v .

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_UNIX).zip

tidy:
	$(GOCMD) mod tidy


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v .
zip-project:
	zip $(BINARY_UNIX).zip $(BINARY_UNIX)
run-localstack:
	docker-compose up --build -d
deploy-localstack:
	aws lambda create-function --function-name gofunction --handler $(BINARY_UNIX) --runtime go1.x --role create-role --zip-file fileb://$(BINARY_UNIX).zip --endpoint-url http://localhost:4566 --region us-east-1 
lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.48-alpine golangci-lint run
clean-localstack:
	docker-compose down --volumes
test:
	go test -race ./...