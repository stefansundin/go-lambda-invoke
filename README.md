This program lets you easily invoke a Go Lambda function locally and send it an event.

How to use:
```
go get -u github.com/stefansundin/go-lambda-invoke

go build
cat event.json | go-lambda-invoke mylambdafunction

# or:
cat event.json | go-lambda-invoke go run mylambdafunction.go
```

See [example](example) for a quick example.

The response from the Lambda function is written to stdout. Both stdout and stderr from the Lambda function is written to stderr.

You may also be interested in [go-lambda-gateway](https://github.com/stefansundin/go-lambda-gateway).
