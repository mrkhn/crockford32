//Copyright 2019 Mark A. Hahn. All rights reserved.
//Use of this source code is governed by an MIT-style
//license that can be found in the LICENSE file.

package crockford32

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type results struct {
		n   int64
		err error
	}

	var tests = []struct {
		args     string
		expected results
	}{
		{args: "", expected: results{err: errDecodeEmptyString("")}},
		{args: "0", expected: results{n: 0}},
		{args: "01", expected: results{n: 1}},
		{args: "10", expected: results{n: 32}},
		{args: "A", expected: results{n: 10}},
		{args: "a", expected: results{n: 10}},
		{args: "*", expected: results{err: errDecodeInvalidRunes("*", "*")}},
		{args: "1!2@3#", expected: results{err: errDecodeInvalidRunes("1!2@3#", "!@#")}},
		{args: "7ZZZZZZZZZZZZ", expected: results{n: math.MaxInt64}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.args), func(t *testing.T) {
			var n, err = Decode(test.args)
			var actual = results{n, err}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf(`Decode(%v) = %d, "%v"; expected %d, "%v"`, test.args, actual.n, actual.err, test.expected.n, test.expected.err)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode("7ZZZZZZZZZZZZ")
	}
}

func TestEncode(t *testing.T) {
	type results struct {
		s   string
		err error
	}

	var tests = []struct {
		args     int64
		expected results
	}{
		{args: 0, expected: results{s: "0", err: nil}},
		{args: -1, expected: results{s: "", err: errEncodeNegativeInt(-1)}},
		{args: 10, expected: results{s: "A", err: nil}},
		{args: 32, expected: results{s: "10", err: nil}},
		{args: math.MaxInt64, expected: results{s: "7ZZZZZZZZZZZZ", err: nil}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.args), func(t *testing.T) {
			var s, err = Encode(test.args)
			var actual = results{s, err}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Encode(%v) = %v, expected %v", test.args, actual, test.expected)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(math.MaxInt64)
	}
}
