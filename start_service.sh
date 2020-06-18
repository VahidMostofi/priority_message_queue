#!/bin/bash
cd service
rm -f main.out
go build  -o main.out .
export API_PORT=7790
export AMQP_URL=amqp://guest:guest@localhost:5672/
export ROLE=SERVICE
export SOURCE_QUEUE=QUEUE_A
export TARGET_QUEUE=QUEUE_B
export WORK_WEIGHT=2
export DEBUG=FALSE
./main.out
