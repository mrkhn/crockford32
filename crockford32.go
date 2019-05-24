//Copyright 2019 Mark A. Hahn. All rights reserved.
//Use of this source code is governed by an MIT-style
//license that can be found in the LICENSE file.

//Package crockford32 is an implementation of Douglas Crockford's
//Base32 specification (https://www.crockford.com/base32.html)
package crockford32

import (
	"fmt"
)

var (
	//encoding is a list of upper case characters used in encoding
	encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	//decoding is an array of characters and the mapped int64 value, used for decoding
	decoding = [128]int64{
		//Singe digit integers runes map to the equivalent integer value
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		//Mappings according to Crockford32 specification ('L' and 'I' map to 1; 'O' maps to 0; 'U' does not map)
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'J': 18, 'K': 19, 'L': 1, 'M': 20, 'N': 21, 'O': 0, 'P': 22, 'Q': 23, 'R': 24, 'S': 25, 'T': 26, 'V': 27, 'W': 28, 'X': 29, 'Y': 30, 'Z': 31,
		//Lower case maps the same as upper case
		'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15, 'g': 16, 'h': 17, 'j': 18, 'k': 19, 'l': 1, 'm': 20, 'n': 21, 'o': 0, 'p': 22, 'q': 23, 'r': 24, 's': 25, 't': 26, 'v': 27, 'w': 28, 'x': 29, 'y': 30, 'z': 31,
	}

	//powers is precalculated powers of 32^i, where i is the index of the slice record
	powers = []int64{
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

//Encode takes an int64 and returns the equivalent Crockford32 encoded string
func Encode(input int64) (string, error) {
	if input == 0 {
		return "0", nil
	}

	if input > 0 {
		count := 0

		runes := make([]byte, 13)
		for ; input > 0; input = input / 32 {
			runes[count] = encoding[input%32]
			count++
		}
		runes = runes[0:count]

		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		return string(runes), nil
	}

	return "", errEncodeNegativeInt(input)
}

//Decode takes a Crockford32 encoded string and returns the equivalent uint64
func Decode(input string) (output int64, err error) {
	if input == "" {
		return 0, errDecodeEmptyString(input)
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
		return 0, errDecodeInvalidRunes(input, invalid[0:count])
	}
	return output, nil
}

func errEncodeNegativeInt(input int64) error {
	return fmt.Errorf(`crockford32.Encode(%d): n must be greater than or equal to zero`, input)
}

func errDecodeEmptyString(input string) error {
	return fmt.Errorf(`crockford32.Decode("%s"): empty string`, input)
}

func errDecodeInvalidRunes(input string, invalid []rune) error {
	return fmt.Errorf(`crockford32.Decode("%s"): contains invalid runes %v`, input, string(invalid))
}
