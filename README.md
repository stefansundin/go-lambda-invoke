This program lets you easily invoke a Go AWS Lambda function locally and send it an event.

Install by running:

```shell
go install github.com/stefansundin/go-lambda-invoke@latest
```

Example usage:

```shell
$ cat event.json | go-lambda-invoke ./mylambdafunction

# or:
$ cat event.json | go-lambda-invoke go run mylambdafunction.go
```

See [example](example) for a quick example.

The response from the Lambda function is written to stdout. Both stdout and stderr from the Lambda function are written to stderr.

You may also be interested in [go-lambda-gateway](https://github.com/stefansundin/go-lambda-gateway).
