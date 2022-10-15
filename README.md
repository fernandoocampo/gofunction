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
    "CodeSize": 7178985,
    "Description": "",
    "Timeout": 3,
    "LastModified": "2022-10-15T16:13:11.484+0000",
    "CodeSha256": "W0zJ2JZdQ8Nm1bxBI6qZPoZpbBxk8pk+1mHOkD6n3aQ=",
    "Version": "$LATEST",
    "VpcConfig": {},
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "13477b74-b32e-4df8-a232-4ec567c94f5a",
    "State": "Active",
    "LastUpdateStatus": "Successful",
    "PackageType": "Zip",
    "Architectures": [
        "x86_64"
    ]
}
```
