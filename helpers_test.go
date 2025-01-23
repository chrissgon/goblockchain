package main

import "testing"

func TestNewStringSHA256(t *testing.T) {
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	have := NewStringSHA256("")

	if have != expected {
		t.Fatalf("should return the hash %v, but got %v", expected, have)
	}
}
