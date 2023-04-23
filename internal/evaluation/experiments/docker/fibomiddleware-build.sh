#!/bin/bash

kind="quic"
echo "Building" $kind

echo "Building namingserver"
#docker build -f evaluation/experiments/docker/fibomiddleware-namingserver.Dockerfile -t midarch/fibomiddleware:1.0.2-namingserver-$kind .
echo
echo

echo "Building server"
#docker build -f evaluation/experiments/docker/fibomiddleware-server.Dockerfile -t midarch/fibomiddleware:1.0.2-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fibomiddleware-client.Dockerfile -t midarch/fibomiddleware:1.0.3-client-$kind .