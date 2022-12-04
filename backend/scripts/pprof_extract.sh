#!/bin/sh

for i in allocs block goroutine heap mutex threadcreate profile trace; do
        echo "Pulling /debug/pprof/$i"
        curl -s http://127.0.0.1:13307/debug/pprof/$i > ./dir/$i
done
