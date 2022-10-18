#!/usr/bin/env bash

# Build application
docker run -w /app -e GOARCH=$(uname -m) -e GOOS=darwin -v $(pwd):/app golang:1.19.2 /usr/local/go/bin/go build -o gopld

# Make executable
chmod +x gopld

# Move to PATH
sudo mv gopld /usr/local/bin/gopld

# Assert pld is in PATH
which gopld
