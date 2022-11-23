# Regex to struct library

[![Go version](https://img.shields.io/github/go-mod/go-version/fclairamb/restruct)](https://golang.org/doc/devel/release.html)
[![Release](https://img.shields.io/github/v/release/fclairamb/restruct)](https://github.com/fclairamb/restruct/releases/latest)
[![Build](https://github.com/fclairamb/restruct/workflows/Build/badge.svg)](https://github.com/fclairamb/restruct/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/fclairamb/restruct/branch/main/graph/badge.svg?token=y1vcrxbXfv)](https://codecov.io/gh/fclairamb/restruct)<!--- [![gocover.io](https://gocover.io/_badge/github.com/fclairamb/restruct)](https://gocover.io/github.com/fclairamb/restruct) -->
[![Go Report Card](https://goreportcard.com/badge/fclairamb/restruct)](https://goreportcard.com/report/fclairamb/restruct)
[![GoDoc](https://godoc.org/github.com/fclairamb/restruct?status.svg)](https://godoc.org/github.com/fclairamb/restruct)

## General idea
This is a very simple library that allows you to convert a regex into a struct. It's intended to be used for simple text parsing around 
dummy bots.

The struct shall have a field for each capture group of the regex.

## Usage

This can be tested [on the playground](https://go.dev/play/p/beFzEua9vlE).

```golang
package main

import (
	"fmt"

	r "github.com/fclairamb/restruct"
)

func main() {

	type Human struct {
		Name   string `restruct:"name"` // Specifying the field
		Age    int    // No tag, "age" will be used
		Height *int   // A pointer will be set to nil if the capture group is empty
	}

	rs := &r.Restruct{
		RegexToStructs: []*r.RegexToStruct{
			{
				ID:     "age",
				Regex:  `^(?P<name>\w+) is (?P<age>\d+)( years old)?$`,
				Struct: &Human{},
			},
			{
				ID:     "height",
				Regex:  `^(?P<name>\w+) is (?P<height>\d+) cm tall$`,
				Struct: &Human{},
			},
		},
	}

	for _, input := range []string{"John is 178 cm tall", "John is 42 years old"} {
		fmt.Println("input:", input)
		m, _ := rs.MatchString(input)

		if m == nil {
			fmt.Printf(`No match for "%s"`, input)

			continue
		}

		fmt.Println("  match ID:", m.ID)

		h := m.Struct.(*Human)
		fmt.Printf("  name = %v, age = %v", h.Name, h.Age)

		if h.Height != nil {
			fmt.Printf(", height = %v", *h.Height)
		}

		fmt.Printf("\n")
	}
}
```
This will produce:
```
input: John is 178 cm tall
  match ID: height
  name = John, age = 0, height = 178
input: John is 42 years old
  match ID: age
  name = John, age = 42
```

## Benchmark
It's _not_ fast:
```text
go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/fclairamb/restruct/test
BenchmarkSmallStruct
BenchmarkSmallStruct-8    	 2813796	       406.5 ns/op	     145 B/op	       6 allocs/op
BenchmarkThreeRules-8     	 1869470	       651.4 ns/op	     145 B/op	       6 allocs/op
BenchmarkBiggerStruct-8   	 2122315	       564.2 ns/op	     177 B/op	      10 allocs/op
```