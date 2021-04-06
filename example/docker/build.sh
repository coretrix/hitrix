#!/usr/bin/env bash

rm -f .env

echo COMPOSE_PROJECT_NAME=hitrix >> .env
echo LOCAL_IP=0.0.0.0 >> .env
echo MYSQL_PORT=9004 >> .env
echo REDIS_PORT=9002 >> .env
echo DEFAULT_REDIS_SEARCH=9002 >> .env

docker-compose up -d --build