package main

import (
	"encoding/hex"
	"fmt"
)

func repeatingKeyXOR(plaintext []byte, key []byte) []byte {
	result := make([]byte, len(plaintext))

	for i := range plaintext {
		result[i] = plaintext[i] ^ key[i%len(key)]
	}
	return result
}

func main() {
	plaintext := []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)

	key := []byte("ICE")

	ciphertext := repeatingKeyXOR(plaintext, key)

	fmt.Println(hex.EncodeToString(ciphertext))
}
