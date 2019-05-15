package crockford32

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type results struct {
		id  string
		err error
	}

	var tests = []struct {
		args     int64
		expected results
	}{
		{args: 0, expected: results{"0", nil}},
		{args: -1, expected: results{"", errEncodeNegativeKey(-1)}},
		{args: 10, expected: results{"A", nil}},
		{args: 32, expected: results{"10", nil}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.args), func(t *testing.T) {
			var id, err = Encode(test.args)
			var actual = results{id, err}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Encode(%v) = %v, expected %v", test.args, actual, test.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type results struct {
		key int64
		err error
	}

	var tests = []struct {
		args     string
		expected results
	}{
		{args: "0", expected: results{key: 0}},
		{args: "A", expected: results{key: 10}},
		{args: "10", expected: results{key: 32}},
		{args: "01", expected: results{key: 1}},
		{args: "", expected: results{err: errDecodeEmptyString("")}},
		{args: "*", expected: results{err: errDecodeInvalidRunes("*", "*")}},
		{args: "1!2@3#", expected: results{err: errDecodeInvalidRunes("1!2@3#", "!@#")}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.args), func(t *testing.T) {
			var key, err = Decode(test.args)
			var actual = results{key, err}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf(`Decode(%v) = %d, "%v"; expected %d, "%v"`, test.args, actual.key, actual.err, test.expected.key, test.expected.err)
			}
		})
	}
}
