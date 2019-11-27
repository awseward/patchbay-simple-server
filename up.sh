#!/usr/bin/env bash

set -euo pipefail

TAG=patchbay-simple-server
PORT=9001

docker build -t "$TAG" .
docker run --rm -p "$PORT:$PORT" "$TAG"
