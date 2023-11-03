#!/bin/bash
CGO_ENABLED=1 GOARCH=amd64 go build .
ln -sf ../lib .

