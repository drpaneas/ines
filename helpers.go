package ines

import (
	"encoding/hex"
	"fmt"
	"os"
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

func asciiDump(hex []string) []string {
	var asciiData []string
	for _, data := range hex {
		asciiData = append(asciiData, hexToASCII(data))
	}
	return asciiData
}

func binaryDump(b []byte) []string {
	var binData []string
	for _, data := range b {
		binData = append(binData, fmt.Sprintf("%08b", data))
	}
	return binData
}

func hexdump(b []byte) []string {
	var hexData []string
	for _, data := range b {
		hexData = append(hexData, byteToHex(data))
	}
	return hexData
}

func byteToHex(b byte) string {
	bs := make([]byte, 1)
	bs[0] = b
	return strings.ToUpper(hex.EncodeToString(bs))
}

func byteToASCII(b byte) string {
	str := fmt.Sprintf("%x", b)
	tmp := fmt.Sprintf("%#v", hexToASCII(str))
	return fixFormat(tmp)
}

func hexToASCII(str string) string {
	bs, err := hex.DecodeString(str)
	if err != nil {
		fmt.Printf("Error: Failed to disassemble: %v", err)
		os.Exit(1)
	}
	return string(bs)
}

func hexToInt(hexStr string) int {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int(result)
}

func intToHex(i int) string {
	return strconv.FormatInt(int64(i), 16)
}

func binToInt(binary string) int64 {
	output, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		fmt.Println("Could not convert binary to integer")
		fmt.Println(err)
		os.Exit(1)
	}
	return output
}

func byteToInt(b byte) int {
	return hexToInt(byteToHex(b))
}

// fixFormat removes the first and the last character of string.
// This solution works even for non-unicode.
func fixFormat(s string) string {
	return trimFLastRune(trimFirstRune(s))
}

func trimFirstRune(s string) string {
	for i := range s {
		if i > 0 {
			// The value i is the index in s of the second
			// rune.  Slice to remove the first rune.
			return s[i:]
		}
	}
	// There are 0 or 1 runes in the string.
	return ""
}

func trimFLastRune(s string) string {
	for i := range s {
		if i > 0 {
			return s[:1]
		}
	}
	// There are 0 or 1 runes in the string.
	return ""

}