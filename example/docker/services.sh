#!/usr/bin/env bash

cd example/docker && docker-compose exec services /bin/sh -c "cd example && APP_MODE=local go run server.go"
