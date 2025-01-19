#!/bin/bash
set -e

# Set your image repository
IMAGE_REPO="mjanith/sr-api-gateway"

# Fetch the latest numerical tag version from Docker Hub
LATEST_TAG=$(curl -s "https://hub.docker.com/v2/repositories/$IMAGE_REPO/tags/"  | \
  jq -r '.results | .[].name' | \
  grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' | \
  sort -V | \
  tail -n 1)

# If there's no latest tag, start from 1.0.0
if [ -z "$LATEST_TAG" ]; then
    NEW_TAG="1.0.0"
else
    # Split the latest tag into parts (assuming the format is major.minor.patch)
    IFS='.' read -r MAJOR MINOR PATCH <<<"${LATEST_TAG//[^0-9.]/}"

    # Increment the PATCH version
    MINOR=$((MINOR + 1))

    # Create the new tag
    NEW_TAG="$MAJOR.$MINOR.$PATCH"
fi

# Build the image with the new tag
docker build -f Dockerfile.prod -t $IMAGE_REPO:$NEW_TAG .

# Push the new image to Docker Hub
docker push $IMAGE_REPO:$NEW_TAG

echo "Pushed image with tag: $NEW_TAG"

