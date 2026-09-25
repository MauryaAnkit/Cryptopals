package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	input1 := "1c0111001f010100061a024b53535009181c"
	input2 := "686974207468652062756c6c277320657965"
	expectedOutput := "746865206b696420646f6e277420706c6179"

	a, err := hex.DecodeString(input1)

	if err != nil {
		panic(err)
	}

	b, err := hex.DecodeString(input2)
	if err != nil {
		panic(err)

	}

	result := fixedXOR(a, b)
	if hex.EncodeToString(result) == expectedOutput {
		fmt.Println("resultMatched")
	}
	fmt.Print(hex.EncodeToString(result))
}

func fixedXOR(a, b []byte) []byte {

	if len(a) != len(b) {
		panic("buffers must be equal length")
	}

	result := make([]byte, len(a))

	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
