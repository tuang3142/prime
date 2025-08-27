package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testInput  = "tests/test.csv"
	testOutput = "tests/test.dat"
)

// TODO: do it
func Test_Encoder(t *testing.T) {
	e := Encoder.New()

	encoded, err := e.Encode(testInput)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	if diff := cmp.Diff(string(file1), string(file2)); diff != "" {
		t.Errorf("Encoded output differs (-want +got):\n%s", diff)
	}
}
