package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

const blockSize = 16

func hasRepeatedBlock(ciphertext []byte) bool {
	seen := make(map[string]bool)

	for i := 0; i < len(ciphertext); i += blockSize {
		// Make sure we have a complete 16-byte block.
		if i+blockSize > len(ciphertext) {
			break
		}

		block := ciphertext[i : i+blockSize]
		key := string(block)

		if seen[key] {
			return true
		}

		seen[key] = true
	}

	return false
}

func main() {
	file, err := os.Open("./set1/challenge08/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()

		// Convert hex string to raw bytes.
		ciphertext, err := hex.DecodeString(line)
		if err != nil {
			fmt.Printf("Line %d: invalid hex: %v\n", lineNumber, err)
			continue
		}

		if hasRepeatedBlock(ciphertext) {
			fmt.Printf("ECB detected on line %d:\n", lineNumber)
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
