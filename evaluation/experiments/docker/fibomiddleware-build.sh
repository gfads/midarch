#!/bin/bash

kind="quic"
echo "Building" $kind

echo "Building namingserver"
#docker build -f evaluation/experiments/docker/fibomiddleware-namingserver.Dockerfile -t midarch/fibomiddleware:namingserver-$kind .
echo
echo

echo "Building server"
docker build -f evaluation/experiments/docker/fibomiddleware-server.Dockerfile -t midarch/fibomiddleware:server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fibomiddleware-client.Dockerfile -t midarch/fibomiddleware:client-$kind .