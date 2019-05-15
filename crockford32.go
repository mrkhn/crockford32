package crockford32

import (
	"fmt"
	"math"
	"strings"
)

var (
	//All characters are encoded in upper case
	encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	decoding = [128]int64{
		//Singe digit integers runes map to the equivalent integer value
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		//Mappings according to Crockford32 specification ('L' and 'I' map to 1; 'O' maps to 0; 'U' does not map)
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'J': 18, 'K': 19, 'L': 1, 'M': 20, 'N': 21, 'O': 0, 'P': 22, 'Q': 23, 'R': 24, 'S': 25, 'T': 26, 'V': 27, 'W': 28, 'X': 29, 'Y': 30, 'Z': 31,
		//Lower case maps the same as upper case
		'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15, 'g': 16, 'h': 17, 'j': 18, 'k': 19, 'l': 1, 'm': 20, 'n': 21, 'o': 0, 'p': 22, 'q': 23, 'r': 24, 's': 25, 't': 26, 'v': 27, 'w': 28, 'x': 29, 'y': 30, 'z': 31,
	}
)

//Encode takes a key and returns the equivalent Crockford32 encoded identifier string
func Encode(key int64) (id string, err error) {
	if key < 0 {
		return id, errEncodeNegativeKey(key)
	}

	if key == 0 {
		return "0", nil
	}

	for i := key; i > 0; i /= 32 {
		id = string(encoding[i%32]) + id
	}

	return id, nil
}

//Decode takes a Crockford32 encoded identifier string and returns the equivalent uint64
func Decode(id string) (key int64, err error) {
	var invalid string
	if id == "" {
		return 0, errDecodeEmptyString(id)
	}

	for i, r := range id {
		if strings.ContainsRune(encoding, r) {
			key += decoding[r] * int64(math.Pow(float64(32), float64(len(id)-i-1)))
		} else {
			invalid += string(r)
		}
	}
	if len(invalid) > 0 {
		return 0, errDecodeInvalidRunes(id, invalid)
	}
	return key, nil
}

func errEncodeNegativeKey(key int64) error {
	return fmt.Errorf(`crockford32.Encode(%d): key must be greater than or equal to zero`, key)
}

func errDecodeEmptyString(id string) error {
	return fmt.Errorf(`crockford32.Decode("%s"): empty string`, id)
}

func errDecodeInvalidRunes(id, runes string) error {
	return fmt.Errorf(`crockford32.Decode("%s"): contains invalid runes [%s]`, id, runes)
}
