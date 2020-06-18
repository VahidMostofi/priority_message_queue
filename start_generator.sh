#!/bin/bash
cd service
rm -f main.out
go build  -o main.out .
export API_PORT=7789
export AMQP_URL=amqp://guest:guest@localhost:5672/
export ROLE=GENERATOR
export TARGET_QUEUE=QUEUE_A
export DURATION=15
export RATE=10
./main.out
