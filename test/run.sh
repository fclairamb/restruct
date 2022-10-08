#!/bin/sh -ex


go test -run=^$ -bench=. -v -benchmem >benchmark.txt

go test -parallel 20 -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=../... ../...
go tool cover -html=coverage.txt

for test in BenchmarkSmallStruct BenchmarkBiggerStruct BenchmarkThreeRules BenchmarkLoadAndExec
do
  go test -bench=$test -run=^$ -cpuprofile $test.pprof -benchtime=5s
  go tool pprof -svg $test.pprof > $test.svg
done
