#!/bin/bash

# Create resources directory if it doesn't exist
mkdir -p resources

# Run the tests
cd /root/project
go test -v ./pkg/tests/
