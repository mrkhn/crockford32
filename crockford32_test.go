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

// TestDecode tests the Decode function by providing a set of inputs and expected outputs.
// It checks if the actual output and error match the expected output and error.
// If they don't match, it reports an error.
// The test cases cover various scenarios such as empty input, valid input, invalid input, and maximum value input.
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

// BenchmarkDecode benchmarks the Decode function with different inputs.
// The first benchmark tests the function with an empty string.
// The second benchmark tests the function with the maximum length string.
// The third benchmark tests the function with a string containing invalid runes.
// The fourth benchmark tests the function with a random input generated using the Encode function.
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

// TestEncode tests the Encode function which encodes an int64 to a base32 string.
// It tests various input values and their expected output values.
// If the actual output or error does not match the expected output or error, the test fails.
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

// BenchmarkEncode benchmarks the Encode function with different input values.
// It tests the performance of encoding the maximum int64 value, zero, negative one, and random int64 values.
// The function uses the testing.B type to run the benchmarks and measure the performance.
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
