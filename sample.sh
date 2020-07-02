#!/bin/sh
# run ./gke-template with parallel
go tool pprof -http=":8081" http://localhost:6060/debug/pprof/profile

open 'http://localhost:6060/debug/pprof/profile'
