// Copyright 2025 Mark A. Hahn. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

// Package crockford32 provides functions to convert base 10 integer values to base 32 strings using Douglas Crockford's Base32 character set, for use as human readable identifiers.
// Crockford32 (https://www.crockford.com/base32.html) is a variant of Base32 encoding that uses a subset of the alphabet to avoid confusion between visually similar characters.
// The 32 character set excludes 'I', 'L', 'O', and 'U', which are easily confused with '1', '1', '0', and 'V', respectively.
package crockford32

import (
	"fmt"
)

var (
	// encoding is a list of upper case characters used for string conversion
	encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// decoding is an array of characters and the mapped uint64 value, used for parsing
	decoding = [128]uint64{
		// Single digit integer runes map to the equivalent integer value
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		// Mappings according to Crockford32 specification ('L' and 'I' map to 1; 'O' maps to 0; 'U' does not map)
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'J': 18, 'K': 19, 'L': 1, 'M': 20, 'N': 21, 'O': 0, 'P': 22, 'Q': 23, 'R': 24, 'S': 25, 'T': 26, 'V': 27, 'W': 28, 'X': 29, 'Y': 30, 'Z': 31,
		// Lower case characters map the same as upper case
		'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15, 'g': 16, 'h': 17, 'j': 18, 'k': 19, 'l': 1, 'm': 20, 'n': 21, 'o': 0, 'p': 22, 'q': 23, 'r': 24, 's': 25, 't': 26, 'v': 27, 'w': 28, 'x': 29, 'y': 30, 'z': 31,
	}

	// powers is precalculated powers of 32^i, where i is the index of the slice record
	powers = []uint64{
		0:  1,
		1:  32,
		2:  1024,
		3:  32768,
		4:  1048576,
		5:  33554432,
		6:  1073741824,
		7:  34359738368,
		8:  1099511627776,
		9:  35184372088832,
		10: 1125899906842624,
		11: 36028797018963968,
		12: 1152921504606846976,
	}
)

// Format converts a base 10 uint64 in to a base 32 number using the Crockford Base32 character.
func Format(input uint64) string {
	if input == 0 {
		return "0"
	}

	count := 0

	runes := make([]byte, 13)

	for ; input > 0; input >>= 5 {
		runes[count] = encoding[input&31]
		count++
	}
	runes = runes[0:count]

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// Parse converts a base 32 number using the Crockford Base32 character set into a base 10 uint64.
// It returns the a uint64 and an error if the input string is empty or contains invalid runes.
func Parse(input string) (output uint64, err error) {
	if input == "" {
		return 0, new(ErrEmptyString)
	}

	if len(input) > len(powers) {
		return 0, ErrInputTooLong{input}
	}

	invalid := make([]rune, len(input))
	count := 0

	for i, r := range input {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z' && r != 'u') || (r >= 'A' && r <= 'Z' && r != 'U') {
			output += decoding[r] * powers[len(input)-i-1]
		} else {
			invalid[count] = r
			count++
		}
	}

	if count > 0 {
		return 0, ErrInvalidRunes{input, string(invalid[0:count])}
	}
	return output, nil
}

// ErrInvalidRunes represents an error that occurs when the input string for decoding contains invalid runes.
type ErrInvalidRunes struct {
	input string
	runes string
}

// Error returns a string representation of the ErrInvalidRunes error.
// It formats the error message to include the input string and the invalid runes found during decoding.
func (e ErrInvalidRunes) Error() string {
	input := e.input
	if len(input) > 32 {
		input = input[:32] + "..."
	}

	runes := e.runes
	if len(runes) > 32 {
		runes = runes[:32] + "..."
	}

	return fmt.Sprintf(`crockford32.Parse(%q): contains invalid runes %q`, input, runes)
}

// ErrEmptyString is an error type that is returned when an empty string is passed to the Crockford32 encoding function.
type ErrEmptyString struct{}

// Error returns the error message for an empty input string.
func (e ErrEmptyString) Error() string {
	return `crockford32.Parse(""): empty string`
}

// ErrInputTooLong is an error type that is returned when a string longer than 13 characters is passed to the Crockford32 encoding function.
type ErrInputTooLong struct {
	input string
}

// Error returns the error message for an input string that is too long.
func (e ErrInputTooLong) Error() string {
	input := e.input
	if len(input) > 32 {
		input = input[:32] + "..."
	}

	return fmt.Sprintf(`crockford32.Parse(%q): input too long`, input)
}
