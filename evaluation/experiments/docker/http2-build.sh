#!/bin/bash

kind="http2"
echo "Building" $kind

echo "Building server"
#docker build -f evaluation/experiments/docker/http2-server.Dockerfile -t midarch/fibomiddleware:1.0.2-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/http2-client.Dockerfile -t midarch/fibomiddleware:1.0.2-client-$kind .