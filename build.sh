#!/bin/sh

set -eu

TAG=$(git log --pretty=%h -n 1)
REPOSITORY=satococoa
NAME=prbot
IMAGE=${REPOSITORY}/${NAME}:${TAG}

docker build . -t $IMAGE
docker push $IMAGE
