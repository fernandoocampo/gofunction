# gofunction
a lambda function in go

## References

* I am using this library [aws-lambda-go](https://github.com/aws/aws-lambda-go/blob/main/events/README_SQS.md).

## Docker compose

* Start dynamodb and create table

```shell
docker-compose up --build
```

* Stop dynamodb

```shell
docker-compose down
```

## DynamoDB

* Scan 

```shell
aws dynamodb scan \
--table-name audit_fruits \
--endpoint-url http://localhost:4566 \
--region us-east-1
```

## Deploy lambda on localstack

```sh
aws lambda create-function --function-name gofunction --handler gofunction-amd64-linux --runtime go1.x --role create-role --zip-file fileb://gofunction-amd64-linux.zip --endpoint-url http://localhost:4566 --region us-east-1

{
    "FunctionName": "gofunction",
    "FunctionArn": "arn:aws:lambda:us-east-1:000000000000:function:gofunction",
    "Runtime": "go1.x",
    "Role": "create-role",
    "Handler": "gofunction-amd64-linux",
    "CodeSize": 7178986,
    "Description": "",
    "Timeout": 3,
    "LastModified": "2022-10-15T16:48:39.002+0000",
    "CodeSha256": "jXcz8Yew9bes6F4d4+fGIqfKQqi5/K64yK9CUpOAam4=",
    "Version": "$LATEST",
    "VpcConfig": {},
    "Environment": {
        "Variables": {
            "CLOUD_ENDPOINT_URL": "http://localhost:4566",
            "CLOUD_REGION": "us-east-1"
        }
    },
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "2a06cda4-f06b-4f84-920f-88013bba5415",
    "State": "Active",
    "LastUpdateStatus": "Successful",
    "PackageType": "Zip",
    "Architectures": [
        "x86_64"
    ]
}
```

## Create sqs

```sh
aws sqs create-queue --queue-name fruits-queue --endpoint-url http://localhost:4566 --region us-east-1

{
    "QueueUrl": "http://localhost:4566/000000000000/fruits-queue"
}
```

* check the queue attributes 

```sh
aws sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/fruits-queue --attribute-names All --endpoint-url http://localhost:4566 --region us-east-1

{
    "Attributes": {
        "ApproximateNumberOfMessages": "0",
        "ApproximateNumberOfMessagesNotVisible": "0",
        "ApproximateNumberOfMessagesDelayed": "0",
        "CreatedTimestamp": "1665852576",
        "DelaySeconds": "0",
        "LastModifiedTimestamp": "1665852576",
        "MaximumMessageSize": "262144",
        "MessageRetentionPeriod": "345600",
        "QueueArn": "arn:aws:sqs:us-east-1:000000000000:fruits-queue",
        "ReceiveMessageWaitTimeSeconds": "0",
        "VisibilityTimeout": "30"
    }
}
```

## Set up event source

We are linking here the queue and the function.

```sh
aws lambda create-event-source-mapping --function-name gofunction --batch-size 5 --maximum-batching-window-in-seconds 60 --event-source-arn arn:aws:sqs:us-east-1:000000000000:fruits-queue --endpoint-url http://localhost:4566 --region us-east-1

{
    "UUID": "c60875e6-fc24-4aab-ba80-fa60968f6104",
    "StartingPosition": "LATEST",
    "BatchSize": 5,
    "ParallelizationFactor": 1,
    "EventSourceArn": "arn:aws:sqs:us-east-1:000000000000:fruits-queue",
    "FunctionArn": "arn:aws:lambda:us-east-1:000000000000:function:gofunction",
    "LastModified": "2022-10-15T18:51:19.497995+02:00",
    "LastProcessingResult": "OK",
    "State": "Enabled",
    "StateTransitionReason": "User action",
    "Topics": [],
    "MaximumRetryAttempts": -1
}
```

## Let's test the function.

Signaling `fruits-queue`.

```sh
aws sqs send-message --queue-url http://localhost:4566/000000000000/fruits-queue --message-body '{"source_id": "1d952b94-a5db-4d63-a500-b486dd96e8b2","name": "lemon","variety": "lima","price": 2.50}' --endpoint-url http://localhost:4566 --region us-east-1

{
    "MD5OfMessageBody": "66eaae934773968fed7241a3629086be",
    "MessageId": "d7f15ca6-aef0-4bae-a5e9-ce676c641943"
}
```

## Scan fruits table

```shell
aws dynamodb scan \
--table-name audit_fruits \
--endpoint-url http://localhost:4566 \
--region us-east-1
```


## The same but using Makefile

```shell
make clean
make clean-localstack
make all
make run-localstack
make deploy-localstack
make create-queue
make event-source-mapping
make signal-fruit
make scan-audit-fruits
make clean-localstack
make clean
```

