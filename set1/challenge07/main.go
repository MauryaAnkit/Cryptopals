package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("./set1/challenge07/input.txt")
	if err != nil {
		panic(err)
	}

	// Base64 decode
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		panic(err)
	}

	// AES-128 with "YELLOW SUBMARINE"
	key := []byte("YELLOW SUBMARINE")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// AES block size = 16 bytes
	plaintext := make([]byte, len(ciphertext))

	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(
			plaintext[i:i+aes.BlockSize],
			ciphertext[i:i+aes.BlockSize],
		)
	}

	fmt.Println(string(plaintext))
}
