#!/bin/bash
set -e

# Use provided environment variables or defaults
IGNITION_HOST=${IGNITION_HOST:-"http://localhost:8088"}

echo "Starting Ignition via Docker Compose..."
docker compose up -d

echo "Waiting for Gateway to be ready at $IGNITION_HOST..."
# Wait for the gateway to respond with gwinfo
timeout 300s bash -c "until curl -s -f $IGNITION_HOST/system/gwinfo; do echo 'Still waiting...'; sleep 10; done"

echo "Ignition Gateway is up!"
