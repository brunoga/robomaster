#!/bin/bash

# Compile example program as an amd64 executable.
CGO_ENABLED=1 GOARCH=amd64 go build .

# Make sure we have a link to the required library so the code will find it.
ln -sf ../wrapper/lib .

