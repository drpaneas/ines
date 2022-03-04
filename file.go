package ines

import (
	"fmt"
	"io/ioutil"
)

// Read reads all the content of the given file path
// and returns it as byte buffer.
func Read(path string) ([]byte, error) {
	var wrapErr error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		wrapErr = fmt.Errorf("failed to read the content of the file %v - Error: %w", path, err)
	}

	return content, wrapErr
}

// Write writes the given byte buffer into disk
// to the given file path.
func Write(path string, buf []byte) error {
	var wrapErr error

	// nolint: gomnd
	if err := ioutil.WriteFile(path, buf, 0o600); err != nil {
		wrapErr = fmt.Errorf("failed to write the content of the buffer to the disk: %v - Error: %w", path, err)
	}

	return wrapErr
}
