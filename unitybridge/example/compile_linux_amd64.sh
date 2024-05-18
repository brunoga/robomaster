#!/bin/bash

# Compile DLL Host as a Windows executable.
GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build ../wrapper/internal/implementations/wine/dllhost

# Compile example program as a Linux executable.
go build .

# Make sure we have a link to the required library so the code will find it.
ln -sf ../wrapper/lib .
