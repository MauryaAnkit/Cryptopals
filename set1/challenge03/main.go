package main

import (
	"encoding/hex"
	"fmt"
	"unicode"
)

// scoreEnglish gives a score to a byte slice.
// Higher score = more likely to be English text.
func scoreEnglish(data []byte) int {
	score := 0

	for _, b := range data {
		switch {
		// space
		case b == ' ':
			score += 5
		// Common English letters.
		case b == 'e' || b == 'E':
			score += 4
		case b == 't' || b == 'T':
			score += 4
		case b == 'a' || b == 'A':
			score += 3
		case b == 'o' || b == 'O':
			score += 3
		case b == 'i' || b == 'I':
			score += 3
		case b == 'n' || b == 'N':
			score += 3
		case b == 's' || b == 'S':
			score += 2
		case b == 'h' || b == 'H':
			score += 2
		case b == 'r' || b == 'R':
			score += 2

		case unicode.IsPrint(rune(b)):
			score += 1

		default:
			score -= 10
		}
	}

	return score
}

// singleByteXOR decrypts ciphertext using a single-byte key.
func singleByteXOR(ciphertext []byte, key byte) []byte {
	result := make([]byte, len(ciphertext))

	for i := range ciphertext {
		result[i] = ciphertext[i] ^ key
	}

	return result
}

func main() {
	input := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	ciphertext, err := hex.DecodeString(input)

	if err != nil {
		panic(err)
	}

	var bestKey byte
	bestScore := -1 << 30
	var bestPlaintext []byte

	for key := 0; key <= 255; key++ {
		plaintext := singleByteXOR(ciphertext, byte(key))

		score := scoreEnglish(plaintext)

		if score > bestScore {
			bestScore = score
			bestKey = byte(key)
			bestPlaintext = plaintext
		}
	}

	fmt.Printf("Key:        0x%02x\n", bestKey)
	fmt.Printf("Plaintext:  %s\n", bestPlaintext)
	fmt.Printf("Score:      %d\n", bestScore)
}
