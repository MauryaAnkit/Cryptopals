package utils

import "fmt"

// FixedXOR XORs two equal-length byte slices.
func FixedXOR(a, b []byte) ([]byte, error) {

	if len(a) != len(b) {
		return nil, fmt.Errorf("inputs must have equal length")
	}

	result := make([]byte, len(a))

	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}

// XORWithKey XORs data with a repeating key.
func XORWithKey(data, key []byte) []byte {
	if len(key) == 0 {
		return nil
	}

	result := make([]byte, len(data))

	for i := range data {
		result[i] = data[i] ^ key[i%len(key)]
	}

	return result
}
