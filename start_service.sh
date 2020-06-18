#!/bin/bash
cd service
go build  -o main .
export API_PORT=7789
export AMQP_URL=amqp://guest:guest@localhost:5672/
export ROLE=SERVICE
export TARGET=TARGET_QUEUE
./main
