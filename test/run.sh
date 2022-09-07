#!/bin/sh -ex
for test in BenchmarkSmallStruct BenchmarkBiggerStruct BenchmarkThreeRules
do
  go test -bench=$test -run=^$ -cpuprofile $test.pprof -benchtime=5s
  go tool pprof -svg $test.pprof > $test.svg
done
go test -run=^$ -bench=. -v -benchmem >benchmark.txt
