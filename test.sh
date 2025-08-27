#!/bin/sh
go build
./mist -i example
open http://localhost:8080 &
./mist -s 8080
