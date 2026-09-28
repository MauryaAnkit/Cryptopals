package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"unicode"
)

type KeySizeScore struct {
	size  int
	score float64
}

func hammingDistance(a, b []byte) int {
	if len(a) != len(b) {
		panic("blocks must have equal length")
	}

	distance := 0

	for i := range a {
		distance += bits.OnesCount8(a[i] ^ b[i])
	}

	return distance
}

func normalizedHammingDistance(data []byte, keySize int) float64 {
	if len(data) < keySize*2 {
		return 999999
	}
	// here 4 is block size it can be different how accurate we want to calculate the average hamming distence
	blocks := make([][]byte, 0, 8)

	for i := 0; i < 8; i++ {
		start := i * keySize
		end := start + keySize

		if end > len(data) {
			break
		}

		blocks = append(blocks, data[start:end])
	}

	if len(blocks) < 2 {
		return 999999
	}

	totalDistance := 0

	for i := 0; i < len(blocks)-1; i++ {
		totalDistance += hammingDistance(blocks[i], blocks[i+1])
	}

	averageDistance := float64(totalDistance) / float64(len(blocks)-1)
	// returing  normalizedHammingDistance
	return averageDistance / float64(keySize)
}

func singleByteXOR(ciphertext []byte, key byte) []byte {
	result := make([]byte, len(ciphertext))

	for i := range ciphertext {
		result[i] = ciphertext[i] ^ key
	}

	return result
}

func scoreEnglish(data []byte) int {
	score := 0

	for _, b := range data {
		switch {
		case b == ' ':
			score += 5
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
		case unicode.IsPrint(rune(b)) || b == '\n' || b == '\r' || b == '\t':
			score += 1
		default:
			score -= 10
		}
	}

	return score
}

func breakSingleByteXOR(ciphertext []byte) (byte, []byte, int) {
	bestScore := -1 << 30
	var bestKey byte
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

	return bestKey, bestPlaintext, bestScore
}

func transpose(ciphertext []byte, keySize int) [][]byte {
	blocks := make([][]byte, keySize)

	for i, b := range ciphertext {
		index := i % keySize
		blocks[index] = append(blocks[index], b)
	}

	return blocks
}

func decryptRepeatingKeyXOR(ciphertext, key []byte) []byte {
	plaintext := make([]byte, len(ciphertext))

	for i := range ciphertext {
		plaintext[i] = ciphertext[i] ^ key[i%len(key)]
	}

	return plaintext
}

func readBase64File(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var encoded []byte

	for scanner.Scan() {
		encoded = append(encoded, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(encoded))
}

func main() {
	ciphertext, err := readBase64File("./set1/challenge06/input.txt")
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------
	// Step 1: Find promising key sizes
	// --------------------------------------------------

	var candidates []KeySizeScore

	for keySize := 2; keySize <= 40; keySize++ {
		score := normalizedHammingDistance(ciphertext, keySize)

		candidates = append(candidates, KeySizeScore{
			size:  keySize,
			score: score,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score < candidates[j].score
	})

	fmt.Println("Top key-size candidates:")

	for i := 0; i < 5; i++ {
		fmt.Printf(
			"KeySize: %d  Normalized Distance: %.4f\n",
			candidates[i].size,
			candidates[i].score,
		)
	}

	// --------------------------------------------------
	// Step 2: Try the top key sizes
	// --------------------------------------------------

	bestOverallScore := -1 << 30
	var bestOverallKey []byte
	var bestOverallPlaintext []byte

	for _, candidate := range candidates[:5] {
		keySize := candidate.size

		blocks := transpose(ciphertext, keySize)

		key := make([]byte, keySize)

		for i, block := range blocks {
			bestKeyByte, _, _ := breakSingleByteXOR(block)
			key[i] = bestKeyByte
		}

		plaintext := decryptRepeatingKeyXOR(ciphertext, key)

		score := scoreEnglish(plaintext)

		if score > bestOverallScore {
			bestOverallScore = score
			bestOverallKey = key
			bestOverallPlaintext = plaintext
		}
	}

	// --------------------------------------------------
	// Step 3: Print result
	// --------------------------------------------------

	fmt.Printf("\nKey: %s\n", bestOverallKey)
	fmt.Printf("Score: %d\n", bestOverallScore)

	fmt.Println("\nPlaintext:")
	fmt.Println(string(bestOverallPlaintext))
}
