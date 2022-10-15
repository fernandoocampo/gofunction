package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fernandoocampo/gofunction/internal/fruits"
)

var fruitsService *fruits.Service

func init() {
	region := os.Getenv("CLOUD_REGION")
	endpointURL := os.Getenv("CLOUD_ENDPOINT_URL")

	newFruitsService, err := fruits.NewService(region, endpointURL)
	if err != nil {
		log.Fatal(err)
	}

	fruitsService = newFruitsService
}

func main() {
	lambda.Start(fruits.MakeHandler(fruitsService))
}
