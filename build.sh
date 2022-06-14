#!/bin/bash
CGO_ENABLED=1 GOOS=linux CC=x86_64-unknown-linux-gnu-gcc CXX=x86_64-unknown-linux-gnu-g++ GOARCH=amd64 go build -ldflags "-X 'github.com/0xJacky/Nginx-UI/server/settings.buildTime=$(date +%s)'" -o nginx-ui-server -v main.go
docker build -t nginx-ui-demo .
docker tag nginx-ui-demo uozi/nginx-ui-demo
docker push uozi/nginx-ui-demo
