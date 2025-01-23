package main

import (
	"errors"
	"testing"
)

func TestBlock_New(t *testing.T) {
	block := NewBlock(0, "", "data").Get()

	if block.Index != 0 {
		t.Fatal("block index should be 0")
	}

	if block.PreviousHash != "" {
		t.Fatal("block should not has previous hash")
	}

	if block.Data != "data" {
		t.Fatal("block data is wrong")
	}
}
func TestBlock_Check(t *testing.T) {
	block := NewBlock(0, "", "").Get()

	err := block.Check()

	if err != nil {
		t.Fatalf("should return nil because block hash is valid, but got %v", err)
	}

	block.Get().Hash = "4fc82b26aecb47d2868c4efbe3581732a3e7cbcc6c2efb32062c08170a05eeb8"

	err = block.Check()

	if !errors.Is(err, ErrInvalidBlockHash) {
		t.Fatalf("should return an invalid block hash error, but got %v", err)
	}
}
