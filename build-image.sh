#!/bin/sh
set -e

IMAGE_NAME="vesicash/auth"
IMAGE_TAG="${1:-latest}"

if [[ ! -f "app.env" ]]
then
    echo "Copying environment file for app ⏳"
    cp app-sample.env app.env
else
    echo "Environment file found 👌"
fi

echo "Building docker image ${IMAGE_TAG} version 🛠️"
docker build -t $IMAGE_NAME:$IMAGE_TAG .