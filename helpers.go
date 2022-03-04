package ines

import (
	"encoding/hex"
	"strconv"
	"strings"
)

func readHighNibbleByte(b byte) byte {
	// To get the high nibble, you shift the value four bits to the right.
	return b >> 4 // nolint: gomnd
}

func readLowNibbleByte(b byte) byte {
	// To get the low nibble, you mask out the lower four bits.
	return b & 15 // nolint: gomnd
}

func mergeNibbles(highNibble byte, lowNibble byte) byte {
	highNibble <<= 4
	return highNibble | lowNibble
}

func byteToHex(b byte) string {
	bs := make([]byte, 1)
	bs[0] = b
	return strings.ToUpper(hex.EncodeToString(bs))
}

// nolint: gomnd
func hexToInt(hexStr string) int {
	// remove 0x suffix if found in the input string
	cleaned := strings.ReplaceAll(hexStr, "0x", "")

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int(result)
}

func byteToInt(b byte) int {
	return hexToInt(byteToHex(b))
}
