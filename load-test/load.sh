#!/bin/bash

URL="http://cusocial.us/feed"  
TOTAL_REQUESTS=10000              
CONCURRENCY=10                    

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
