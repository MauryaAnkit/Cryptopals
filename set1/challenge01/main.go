package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func hexToBase64(input string) (string, error) {
	rawBytes, err := hex.DecodeString(input)
	if err != nil {
		fmt.Print("Error", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(rawBytes), nil

}

func main() {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	expectedOutput := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	output, err := hexToBase64(input)

	if err != nil {
		fmt.Print("Error", err)
		return
	}

	if output == expectedOutput {
		fmt.Println("resultMatched")
	}
	fmt.Println(output)
}
