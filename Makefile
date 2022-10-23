# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=gofunction
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
DOCKER_REPO=fdocampo
DOCKER_CONTAINER=frutal

all: clean build-linux zip-project

build-and-push: build-linux zip-project push-build-s3

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
push-build-s3:
	aws s3 cp ./$(BINARY_UNIX).zip s3://function-bucket --endpoint-url http://localhost:4566 --region us-east-1
run-localstack:
	docker-compose up --build -d
deploy-localstack:
	aws lambda create-function --function-name gofunction --environment "Variables={CLOUD_REGION=us-east-1,CLOUD_ENDPOINT_URL=http://localhost:4566}" --handler $(BINARY_UNIX) --runtime go1.x --role create-role --zip-file fileb://$(BINARY_UNIX).zip --endpoint-url http://localhost:4566 --region us-east-1 
create-queue:
	aws sqs create-queue --queue-name audit-fruits --endpoint-url http://localhost:4566 --region us-east-1
create-table:
	aws dynamodb create-table \
	--table-name audit-fruits \
	--attribute-definitions AttributeName=id,AttributeType=S \
	--key-schema AttributeName=id,KeyType=HASH \
	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
	--endpoint-url http://localhost:4566 --region us-east-1
queue-attributes:
	aws sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/audit-fruits --attribute-names All --endpoint-url http://localhost:4566 --region us-east-1
event-source-mapping:
	aws lambda create-event-source-mapping --function-name gofunction --batch-size 5 --maximum-batching-window-in-seconds 60 --event-source-arn arn:aws:sqs:us-east-1:000000000000:audit-fruits --endpoint-url http://localhost:4566 --region us-east-1
signal-fruit:
	aws sqs send-message --queue-url http://localhost:4566/000000000000/audit-fruits --message-body '{"source_id": "1d952b94-a5db-4d63-a500-b486dd96e8b2","name": "lemon-sqs","variety": "lima-sqs","price": 2.50}' --endpoint-url http://localhost:4566 --region us-east-1
scan-audit-fruits:
	aws dynamodb scan --table-name audit-fruits --endpoint-url http://localhost:4566 --region us-east-1
clean-localstack:
	docker-compose down --volumes
lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.48-alpine golangci-lint run
test:
	go test -race ./...