package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"sync"
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

	// Read the input data
	payload, _ := io.ReadAll(os.Stdin)

	// When it is time to exit, wait for the the Lambda process to be killed
	var wg sync.WaitGroup
	var exitCode int
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
		os.Exit(exitCode)
	}()
	// Run the Lambda function in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.CommandContext(ctx, os.Args[1], os.Args[2:]...)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(),
			"_LAMBDA_SERVER_PORT="+strconv.Itoa(port),
		)
		err := cmd.Run()
		if err == nil {
			// This case should not happen normally
			// If it does, it means that the process exited on its own, which a Lambda function shouldn't do
			// This can happen if you run a non-Lambda program (e.g. /bin/ls)
			fmt.Fprintf(os.Stderr, "Warning: Lambda process finished unexpectedly.\n")
			os.Exit(1)
		}
		// Ignore the error if it was caused by us
		if e, ok := err.(*exec.ExitError); ok && e.Error() == "signal: killed" {
			return
		}
		panic(err)
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
	deadline := time.Now().Add(60 * time.Minute)
	invokeRequest := &messages.InvokeRequest{
		Payload:      payload,
		RequestId:    "0",
		XAmznTraceId: "",
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: int64(deadline.Unix()),
			Nanos:   int64(deadline.Nanosecond()),
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
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "The Lambda function encountered an error:\n")
		fmt.Fprintln(os.Stderr, invokeResponse.Error.Message)
		for _, frame := range invokeResponse.Error.StackTrace {
			fmt.Fprintf(os.Stderr, "\t%s:%d %s\n", frame.Path, frame.Line, frame.Label)
		}
		exitCode = 1
		return
	}

	// Print the output
	fmt.Print(string(invokeResponse.Payload))
	// If the output does not end with a newline, then output one to stderr to make the terminal look nicer
	if invokeResponse.Payload[len(invokeResponse.Payload)-1] != '\n' {
		fmt.Fprint(os.Stderr, "\n")
	}
}
