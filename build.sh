#!/bin/bash

export GOPATH=$(pwd)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo httpProxy
docker build -t spectralogic/sdk_proxy:latest .

# export GOPATH=$PWD; go install httpProxy
# docker run -p 9080:9080 -p 9090:9090 spectralogic/sdk_proxy:latest
# docker tag 2aacb1b2635b spectralogic/sdk_proxy:latest
# docker login <creds>
# docker push spectralogic/sdk_proxy:latest
# curl -X PUT http://localhost:9090/close

