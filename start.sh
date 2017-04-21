#! /bin/bash

cd bin && go build ../src/main.go
cd ../ && CONFIGOR_ENV=`cat .env` ./bin/main

