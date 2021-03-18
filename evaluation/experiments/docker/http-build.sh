#!/bin/bash

kind="https"
echo "Building" $kind

echo "Building server"
docker build -f evaluation/experiments/docker/http-server.Dockerfile -t midarch/fibomiddleware:server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/http-client.Dockerfile -t midarch/fibomiddleware:client-$kind .