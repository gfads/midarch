#!/bin/bash

echo "Building docker monitor"

#echo "Building server"
docker build -f evaluation/experiments/docker/monitorExperiments/dockerExperiments.Dockerfile -t midarch/monitor:1.0.0 .
echo
echo
