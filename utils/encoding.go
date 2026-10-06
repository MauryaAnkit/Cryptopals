package utils

import (
	"encoding/base64"
	"encoding/hex"
)

// HexToBytes converts a hexadecimal string into raw bytes.
func HexToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

// BytesToBase64 converts raw bytes into a Base64-encoded string.
func BytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64ToBytes converts a Base64-encoded string into raw bytes.
func Base64ToBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// BytesToHex converts raw bytes into a hexadecimal string.
func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}
