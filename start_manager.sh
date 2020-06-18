#!/bin/bash 
cd service 
rm -f main.out
go build  -o main.out . 
export AMQP_URL=amqp://guest:guest@localhost:5672/ 
export ROLE=MANAGER
export LOAD_GENERATOR_URL=localhost
export API_PORT=7789
export FINAL_QUEUE=QUEUE_B
export DURATION=15
export RATE=1
export DEBUG=FALSE
./main.out
