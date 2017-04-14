#!/bin/sh

set -e
go build -o test .
./test

#strace -o aplay.log aplay -L
#strace -o test.log ./test

#cat aplay.log | tail -n +78 > aplay2.log
#cat test.log | tail -n +143 > test2.log

