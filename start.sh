#! /bin/bash

#cd bin && GOOS=linux GOARCH=amd64 go build ../src/main.go
cd bin && go build ../src/main.go
cd ../ && CONFIGOR_ENV=`cat .env` ./bin/main

