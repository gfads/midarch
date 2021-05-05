#!/bin/bash

kind="rmq"
echo "Building non gMidArch app" $kind

echo "Building server"
docker build -f evaluation/experiments/docker/fibormq-server.Dockerfile -t midarch/fibormq:1.0.3-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fibormq-client.Dockerfile -t midarch/fibormq:1.0.3-client-$kind .