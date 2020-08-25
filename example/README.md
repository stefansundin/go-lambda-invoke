Run this simple Lambda function with an example event:

```
$ cat event.json | go-lambda-invoke go run toupperlambda.go
Hello world!
{"statusCode":200,"headers":{"Content-Type":"text/plain"},"multiValueHeaders":null,"body":"HELLO WORLD!\n"}
```

The first output line, `Hello world!`, is printed in the Lambda function. This would appear in CloudWatch logs.

The second line is what the Lambda function sent in response to the event. It would display `HELLO WORLD!` in a web browser.
