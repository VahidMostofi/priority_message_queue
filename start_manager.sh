#!/bin/bash 
cd service 
go build  -o main.out . 
export AMQP_URL=amqp://guest:guest@localhost:5672/ 
export ROLE=MANAGER 
./main
