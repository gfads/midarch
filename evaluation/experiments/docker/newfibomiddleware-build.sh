#!/bin/bash

kind="udptcp"
echo "Building" $kind

echo "Building namingserver"
docker build -f evaluation/experiments/docker/newfibomiddleware-namingserver.Dockerfile -t midarch/newfibomiddleware:1.0.4-namingserver-$kind .
echo
echo

echo "Building server"
docker build -f evaluation/experiments/docker/newfibomiddleware-server.Dockerfile -t midarch/newfibomiddleware:1.0.4-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/newfibomiddleware-client.Dockerfile -t midarch/newfibomiddleware:1.0.4-client-$kind .