package main

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file to test
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	tests := []struct {
		filename string
		expected bool
	}{
		{filename: tempFile.Name(), expected: true},
		{filename: "nonexistentfile.lua", expected: false},
	}

	for _, test := range tests {
		result := fileExists(test.filename)
		if result != test.expected {
			t.Errorf("fileExists(%s) = %v; want %v", test.filename, result, test.expected)
		}
	}
}
