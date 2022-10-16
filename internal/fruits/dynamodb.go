package fruits

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const fruitsTable = "audit-fruits"

// Setup contains dynamodb settings.
type Setup struct {
	Region   string
	Endpoint string
}

// DynamoDB defines logic for dynamodb repository.
type DynamoDB struct {
	client *dynamodb.Client
}

var (
	errLoadingAWSConfig = errors.New("unable to load aws config")
	errSavingFruit      = errors.New("unable to save fruit")
)

func newDynamoDBClient(ctx context.Context, setup Setup) (*DynamoDB, error) {
	newDynamodb := new(DynamoDB)

	awsconfig, err := newDynamodb.getConfig(ctx, setup.Region, setup.Endpoint)
	if err != nil {
		log.Fatalln("unable to create dynamodb client", err)

		return nil, errLoadingAWSConfig
	}

	newDynamodb.client = dynamodb.NewFromConfig(awsconfig)

	return newDynamodb, nil
}

func (d *DynamoDB) getConfig(ctx context.Context, region, endpoint string) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}

		return aws.Endpoint{}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx, config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("d", "d", "")),
	)
	if err != nil {
		log.Println("unable to load aws config", err)

		return cfg, errLoadingAWSConfig
	}

	return cfg, nil
}

func (d *DynamoDB) Save(ctx context.Context, fruit Fruit) error {
	data, err := attributevalue.MarshalMap(fruit)
	if err != nil {
		log.Println("unable to marshal new fruit", err)

		return errSavingFruit
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(fruitsTable),
		Item:      data,
	})
	if err != nil {
		log.Println("unable to store fruit", err)

		return errSavingFruit
	}

	return nil
}
