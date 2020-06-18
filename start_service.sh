#!/bin/bash
cd service
rm -f main.out
go build  -o main.out .
export API_PORT=7789
export AMQP_URL=amqp://guest:guest@localhost:5672/
export ROLE=SERVICE
export TARGET_QUEUE=QUEUE_B
./main.out
