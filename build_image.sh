#!/bin/bash
#
# A simple script to build the Docker images.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG}
# Note: TAG is not used from workflows, always generated uniquely
set -ex

DOCKERFILE=$1
CONTEXT_DIR=$(dirname "$DOCKERFILE")
IMAGE=$2

cd $CONTEXT_DIR
TAG=$(git describe --tags --always --dirty)

gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Building ${IMAGE} using gcloud build"
gcloud builds submit --tag=${IMAGE}:${TAG} --project=${GCP_PROJECT} .
echo "Finished building ${IMAGE}:${TAG}"
