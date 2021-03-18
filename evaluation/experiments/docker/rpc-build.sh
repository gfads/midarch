#!/bin/bash

kind="rpc"
echo "Building" $kind

echo "Building server"
docker build -f evaluation/experiments/docker/rpc-server.Dockerfile -t midarch/fibomiddleware:server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/rpc-client.Dockerfile -t midarch/fibomiddleware:client-$kind .