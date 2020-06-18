#!/bin/bash 
cd service 
go build  -o main . 
export AMQP_URL=amqp://guest:guest@localhost:5672/ 
export ROLE=MANAGER 
./main
