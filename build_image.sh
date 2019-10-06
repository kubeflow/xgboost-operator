#!/bin/bash
#
# A simple script to build the Docker images.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG}
set -ex

DOCKERFILE=$1
CONTEXT_DIR=$(dirname "$DOCKERFILE")
IMAGE=$2
TAG=$3

gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

cd $CONTEXT_DIR

echo "Building ${IMAGE} using gcloud build"
gcloud builds submit --tag=${IMAGE}:${TAG} --project=${GCP_PROJECT} .
echo "Finished building ${IMAGE}:${TAG}"
