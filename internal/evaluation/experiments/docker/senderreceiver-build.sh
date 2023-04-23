#!/bin/bash

kind="senderreceiver"
echo "Building" $kind

docker build -f evaluation/experiments/docker/senderreceiver.Dockerfile -t senderreceiver .