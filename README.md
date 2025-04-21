# crockford32
Douglas Crockford's Base32 encoding for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/mrkhn/crockford32.svg)](https://pkg.go.dev/github.com/mrkhn/crockford32)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**crockford32** is a Go library that provides a simple and efficient implementation of the Douglas [Crockford's Base32](https://en.wikipedia.org/wiki/Base32#Crockford's_Base32) encoding scheme. This encoding is designed to be human-readable and unambiguous, making it suitable for various applications where data needs to be represented in a text-friendly format.

## Features
* Format uint64 values in to Crockford Base32 encoded strings.
* Parse Crockford Base32 encoded strings in to uint64 values.
* Invalid runes are returned in an error from Parse() logging and debugging.
* Unit tests, benchmarks and example functions.
* MIT license.

## Installation

```sh
go get github.com/yourusername/crockford32
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/mrkhn/crockford32"
)

func main() {
    //Parsing a valid string in to a uint64 value.
    i, err := Parse("10")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(i)
	// Output: 32

    // Formatting a uint64 it's Crockford Base32 string.
	s := Format(32)
	fmt.Println(s)
	// Output: 10
}
```