This program lets you easily invoke a Lambda function locally and send it an event.

How to use:
```
go get -u github.com/stefansundin/go-lambda-invoke

go build
cat event.json | go-lambda-invoke mylambdafunction
```

See [example](example) for a quick example.

You may also be interested in [go-lambda-gateway](https://github.com/stefansundin/go-lambda-gateway).
