#!/usr/bin/env bash
APP_NAME=ldbc
cd $APP_NAME
GOOS=windows GOARCH=amd64 go build -o ../bin/windows/$APP_NAME.exe
GOOS=darwin GOARCH=amd64 go build -o ../bin/macos/$APP_NAME
GOOS=linux GOARCH=amd64 go build -o ../bin/linux/$APP_NAME