package utils

import "fmt"

// PKCS7Pad adds PKCS#7 padding to data so its length
// becomes a multiple of blockSize.
func PKCS7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 {
		return nil, fmt.Errorf("invalid block size: %d", blockSize)
	}

	paddingLength := blockSize - (len(data) % blockSize)
	padding := byte(paddingLength)

	result := make([]byte, len(data)+paddingLength)
	copy(result, data)

	for i := len(data); i < len(result); i++ {
		result[i] = padding
	}

	return result, nil
}
