//Copyright 2019 Mark A. Hahn. All rights reserved.
//Use of this source code is governed by an MIT-style
//license that can be found in the LICENSE file.

package crockford32

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestDecode(t *testing.T) {
	type results struct {
		output int64
		err    error
	}

	tests := []struct {
		input    string
		expected results
	}{
		{input: "", expected: results{output: 0, err: errDecodeEmptyString("")}},
		{input: "0", expected: results{output: 0, err: nil}},
		{input: "01", expected: results{output: 1, err: nil}},
		{input: "10", expected: results{output: 32, err: nil}},
		{input: "A", expected: results{output: 10, err: nil}},
		{input: "a", expected: results{output: 10, err: nil}},
		{input: "*", expected: results{output: 0, err: errDecodeInvalidRunes("*", []rune{'*'})}},
		{input: "1!2@3#", expected: results{output: 0, err: errDecodeInvalidRunes("1!2@3#", []rune{'!', '@', '#'})}},
		{input: "7ZZZZZZZZZZZZ", expected: results{output: math.MaxInt64}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			output, err := Decode(test.input)
			actual := results{output, err}
			if actual.output != test.expected.output || !errors.Is(actual.err, test.expected.err) {
				t.Errorf(`Decode(%v) = %d, "%v"; expected %d, "%v"`, test.input, actual.output, actual.err, test.expected.output, test.expected.err)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	b.Run("Empty", func(b *testing.B) {
		Decode("")
	})
	b.Run("MaxLength", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Decode("7ZZZZZZZZZZZZ")
		}
	})
	b.Run("InvalidRunes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Decode("!@#$%^&*()?+~`-',.-;}{=_|")
		}
	})
	b.Run("Random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			input, _ := Encode(int64(rand.Uint32()))
			Decode(input)
		}
	})
}

func TestEncode(t *testing.T) {
	type results struct {
		output string
		err    error
	}

	tests := []struct {
		input    int64
		expected results
	}{
		{input: 0, expected: results{output: "0", err: nil}},
		{input: -1, expected: results{output: "", err: errEncodeNegativeInt(-1)}},
		{input: 10, expected: results{output: "A", err: nil}},
		{input: 32, expected: results{output: "10", err: nil}},
		{input: math.MaxInt64, expected: results{output: "7ZZZZZZZZZZZZ", err: nil}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			output, err := Encode(test.input)
			actual := results{output, err}
			if actual != test.expected || !errors.Is(actual.err, test.expected.err) {
				t.Errorf("Encode(%v) = %v, expected %v", test.input, actual, test.expected)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	b.Run("MaxInt64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Encode(math.MaxInt64)
		}
	})

	b.Run("Zero", func(b *testing.B) {
		Encode(0)
	})

	b.Run("Negative", func(b *testing.B) {
		Encode(-1)
	})

	b.Run("Random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Encode(int64(rand.Uint32()))
		}
	})
}
