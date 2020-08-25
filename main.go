package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func getFreeTCPPort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Please supply a path to a compiled Go Lambda function. Provide input on stdin.\n")
		os.Exit(1)
	}

	// Get a free port
	port, err := getFreeTCPPort()
	if err != nil {
		panic(err)
	}
	host := "localhost:" + strconv.Itoa(port)

	// Get the absolute path to the binary
	binary, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Read the input data
	payload, _ := ioutil.ReadAll(os.Stdin)

	// Run the Lambda function in a go routine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		cmd := exec.CommandContext(ctx, binary)
		cmd.Env = append(os.Environ(),
			"_LAMBDA_SERVER_PORT="+strconv.Itoa(port),
		)
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}()

	// Wait for the Lambda function to become ready
	var client *rpc.Client
	for {
		client, err = rpc.Dial("tcp", host)
		if err != nil {
			time.Sleep(time.Millisecond)
			continue
		}
		break
	}
	for {
		var pingResponse messages.PingResponse
		if err := client.Call("Function.Ping", &messages.PingRequest{}, &pingResponse); err != nil {
			time.Sleep(time.Millisecond)
			continue
		}
		break
	}

	// Invoke the Lambda function
	var invokeResponse messages.InvokeResponse
	now := time.Now()
	invokeRequest := &messages.InvokeRequest{
		Payload:      payload,
		RequestId:    "0",
		XAmznTraceId: "",
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: int64(now.Unix()),
			Nanos:   int64(now.Nanosecond()),
		},
		InvokedFunctionArn:    "",
		CognitoIdentityId:     "",
		CognitoIdentityPoolId: "",
		ClientContext:         nil,
	}
	if err = client.Call("Function.Invoke", invokeRequest, &invokeResponse); err != nil {
		panic(err)
	}
	if invokeResponse.Error != nil {
		panic(invokeResponse.Error.Message)
	}

	// Print the output
	fmt.Println(string(invokeResponse.Payload))
}
