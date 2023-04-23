#!/bin/bash

kind="loadPlugin"
echo "Building" $kind

docker build -f evaluation/experiments/docker/loadPlugin.Dockerfile -t loadplugin .