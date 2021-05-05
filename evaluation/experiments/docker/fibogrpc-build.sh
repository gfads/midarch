#!/bin/bash

kind="grpc"
echo "Building non gMidArch app" $kind

echo "Building server"
docker build -f evaluation/experiments/docker/fibogrpc-server.Dockerfile -t midarch/fibogrpc:1.0.3-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fibogrpc-client.Dockerfile -t midarch/fibogrpc:1.0.3-client-$kind .