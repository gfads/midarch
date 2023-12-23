#!/bin/bash

kind="tcp"
echo "Building" $kind

echo "Building namingserver"
docker build -f evaluation/experiments/docker/fibonaccidistributed-namingserver.Dockerfile -t midarch/fibonaccidistributed:1.14.0-namingserver-$kind .
echo
echo

echo "Building server"
docker build -f evaluation/experiments/docker/fibonaccidistributed-server.Dockerfile -t midarch/fibonaccidistributed:1.14.0-server-$kind .
echo
echo

echo "Building client"
docker build -f evaluation/experiments/docker/fibonaccidistributed-client.Dockerfile -t midarch/fibonaccidistributed:1.14.0-client-$kind .