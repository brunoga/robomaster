#/bin/sh

# Under darwin/arm64 we compile as amd64 and run under Rosetta (the unitybridge
# library is only available for darwin/amd64).
GOARCH=amd64 CGO_ENABLED=1 go build


