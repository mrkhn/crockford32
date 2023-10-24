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

// TestParse tests the Parse function by providing a set of inputs and expected outputs.
// It checks if the actual output and error match the expected output and error.
// If they don't match, it reports an error.
// The test cases cover various scenarios such as empty input, valid input, invalid input, and maximum value input.
func TestParse(t *testing.T) {
	type results struct {
		output uint64
		err    error
	}

	tests := []struct {
		input    string
		expected results
	}{
		{input: "", expected: results{output: 0, err: new(ErrEmptyString)}},
		{input: "0", expected: results{output: 0, err: nil}},
		{input: "01", expected: results{output: 1, err: nil}},
		{input: "10", expected: results{output: 32, err: nil}},
		{input: "A", expected: results{output: 10, err: nil}},
		{input: "a", expected: results{output: 10, err: nil}},
		{input: "*", expected: results{output: 0, err: ErrInvalidRunes{"*", "*"}}},
		{input: "1!2@3#", expected: results{output: 0, err: ErrInvalidRunes{"1!2@3#", "!@#"}}},
		{input: "7ZZZZZZZZZZZZ", expected: results{output: math.MaxInt64}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			output, err := Parse(test.input)
			actual := results{output, err}
			if actual.output != test.expected.output || !errors.Is(actual.err, test.expected.err) {
				t.Errorf(`Parse(%v) = %d, "%v"; expected %d, "%v"`, test.input, actual.output, actual.err, test.expected.output, test.expected.err)
			}
		})
	}
}

// BenchmarkParse benchmarks the Parse function with different inputs.
// The first benchmark tests the function with an empty string.
// The second benchmark tests the function with the maximum length string.
// The third benchmark tests the function with a string containing invalid runes.
// The fourth benchmark tests the function with a random input generated using the Format function.
func BenchmarkParse(b *testing.B) {
	b.Run("Empty", func(b *testing.B) {
		Parse("")
	})
	b.Run("MaxLength", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Parse("7ZZZZZZZZZZZZ")
		}
	})
	b.Run("InvalidRunes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Parse("!@#$%^&*()?+~`-',.-;}{=_|")
		}
	})
	b.Run("Random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			input := Format(rand.Uint64())
			Parse(input)
		}
	})
}

func TestFormat(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 10, expected: "A"},
		{input: 32, expected: "10"},
		{input: math.MaxInt64, expected: "7ZZZZZZZZZZZZ"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			actual := Format(test.input)

			if actual != test.expected {
				t.Errorf("Format(%v) = %v, expected %v", test.input, actual, test.expected)
			}
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	b.Run("MaxInt64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Format(math.MaxInt64)
		}
	})

	b.Run("Zero", func(b *testing.B) {
		Format(0)
	})

	b.Run("Random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Format(rand.Uint64())
		}
	})
}
