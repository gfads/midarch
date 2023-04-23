#!/bin/bash

kind="rpc"
echo "Building non gMidArch app" $kind

echo "Building server"
docker build -f evaluation/experiments/docker/fiborpc-server.Dockerfile -t midarch/fiborpc:1.0.3-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fiborpc-client.Dockerfile -t midarch/fiborpc:1.0.3-client-$kind .