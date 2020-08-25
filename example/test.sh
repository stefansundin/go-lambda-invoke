#!/bin/bash -e
go build
cat event.json | go-lambda-invoke toupperlambda
