package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
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
	file, err := os.Open("./set1/challenge04/input.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	bestOverallScore := -1 << 30
	var bestOverallKey byte
	var bestOverallPlaintext []byte

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		ciphertext, err := hex.DecodeString(line)

		if err != nil {
			panic(err)
		}

		var bestLineKey byte
		bestLineScore := -1 << 30
		var bestLinePlaintext []byte

		for key := 0; key <= 255; key++ {
			plaintext := singleByteXOR(ciphertext, byte(key))

			score := scoreEnglish(plaintext)

			if score > bestLineScore {
				bestLineScore = score
				bestLineKey = byte(key)
				bestLinePlaintext = plaintext
			}
		}

		if bestLineScore > bestOverallScore {
			bestOverallScore = bestLineScore
			bestOverallKey = bestLineKey
			bestOverallPlaintext = bestLinePlaintext
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Key:        0x%02x\n", bestOverallKey)
	fmt.Printf("Plaintext:  %s", bestOverallPlaintext)
	fmt.Printf("Score:      %d", bestOverallScore)
}
