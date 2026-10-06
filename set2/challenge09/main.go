package main

import (
	"fmt"
	"log"

	"cryptopals/utils"
)

func main() {
	plaintext := []byte("YELLOW SUBMARINE")
	blockSize := 16

	padded, err := utils.PKCS7Pad(plaintext, blockSize)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original: %s\n", plaintext)
	fmt.Printf("Padded:   %q\n", padded)
	fmt.Printf("Hex:      %x\n", padded)
	fmt.Printf("Length:   %d bytes\n", len(padded))
}
