#!/bin/bash -e
cat event.json | go-lambda-invoke go run toupperlambda.go
