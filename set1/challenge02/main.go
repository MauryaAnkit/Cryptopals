package main

import (
	"cryptopals/utils"
	"fmt"
)

func main() {
	input1 := "1c0111001f010100061a024b53535009181c"
	input2 := "686974207468652062756c6c277320657965"
	expectedOutput := "746865206b696420646f6e277420706c6179"

	a, err := utils.HexToBytes(input1)

	if err != nil {
		panic(err)
	}

	b, err := utils.HexToBytes(input2)
	if err != nil {
		panic(err)

	}

	byteResult, err := utils.FixedXOR(a, b)
	if err != nil {
		panic(err)

	}

	result := utils.BytesToHex(byteResult)
	// result := utils.BytesToHex(fixedXOR(a, b))
	if result == expectedOutput {
		fmt.Println("resultMatched")
	}
	fmt.Print(result)
}
