#!/bin/bash

URL="http://34.160.40.147/feed"  # Replace with your endpoint URL
TOTAL_REQUESTS=10000                # Total number of requests to send
CONCURRENCY=10                    # Number of concurrent requests

# Function to perform concurrent requests
send_requests() {
  for ((i=0; i<$CONCURRENCY; i++)); do
    curl -s -o /dev/null -w "%{http_code}\n" $URL &
  done
  wait
}

echo "Starting load test..."
for ((j=0; j<$TOTAL_REQUESTS; j+=$CONCURRENCY)); do
  send_requests
done

echo "Load test completed."
