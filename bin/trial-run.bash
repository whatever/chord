#!/bin/bash

PORTS=(8081 8082 8083 8084)


clean_up() {
  sleep 1
  kill -9 $PID
  echo "<<<FINISHED>>>"
}

trap clean_up SIGINT

go run run.go -l 8081 &
PID=$$
echo Starting seed server on $PID

sleep 1
seq 8082 8093 | xargs -I{} -P200 go run run.go -seeds=":8081" -l {}
