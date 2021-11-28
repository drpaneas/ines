// Package ines implements a simple library for parsing NES roms based on iNES specification
package ines

import (
	"bytes"
	"encoding/hex"
)

// hasHeader returns true if input starts with 'NES^Z' (Hex equiv: 0x4e 0x45 0x53 0x1a).
func hasHeader(b []byte) bool {
	return bytes.HasPrefix(b, bytes.Trim(hexBytes("4e 45 53 1a"), " "))
}

// isINES2 returns true if the 7th Byte has bit-3 set and bit-2 off.
func isINES2(b []byte) bool {
	return hasBit(b[7],3) && !hasBit(b[7],2)
}

// hasBit returns true if b byte has bit-p set.
func hasBit(b byte, p uint8) bool {
	return b & (1 <<p) > 0
}

func hexBytes(h string) []byte {
	return func(b []byte, _ error) []byte {return b}(hex.DecodeString(h))
}
