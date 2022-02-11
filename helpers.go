package ines

import (
	"encoding/hex"
	"strconv"
	"strings"
)

func readHighNibbleByte(B byte) byte {
	B = B >> 4
	if B > 15 {
		B = 15
	}
	return B
}

func readLowNibbleByte(B byte) byte {
	B = B << 4
	B = B >> 4
	if B > 15 {
		B = 15
	}
	return B
}

func mergeNibbles(highNibble byte, lowNibble byte) byte {
	highNibble = highNibble << 4
	return highNibble | lowNibble
}

func byteToHex(b byte) string {
	bs := make([]byte, 1)
	bs[0] = b
	return strings.ToUpper(hex.EncodeToString(bs))
}

func hexToInt(hexStr string) int {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int(result)
}

func byteToInt(b byte) int {
	return hexToInt(byteToHex(b))
}
