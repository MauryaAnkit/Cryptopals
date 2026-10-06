package main

import (
	"cryptopals/utils"
	"fmt"
)

func main() {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	expectedOutput := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	rawBytes, err := utils.HexToBytes(input)

	if err != nil {
		fmt.Print("Error", err)
		return
	}
	result := utils.BytesToBase64(rawBytes)
	if result == expectedOutput {
		fmt.Println("resultMatched")
	}
	fmt.Println(result)
}
