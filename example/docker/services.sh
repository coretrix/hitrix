#!/usr/bin/env bash

cd example/docker && docker-compose exec services /bin/sh -c "cd example/cmd/$1 && APP_MODE=local go run main.go"
