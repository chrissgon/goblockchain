package main

import (
	"errors"
	"testing"
	"time"
)

const POW_LIMIT = 100

func TestPoW_Mine_CheckHash(t *testing.T) {
	pow := NewPoW(0, 0)

	block := NewBlock(0, "", "")

	powBlock, err := pow.Mine(block)

	if err != nil {
		t.Fatal(err)
	}

	hash := NewStringSHA256(GetBlockData(*block.Get()))

	if powBlock.Get().Hash != hash {
		t.Fatal("should return an valid hash")
	}
}
func TestPoW_Mine_NoDifficulty(t *testing.T) {
	pow := NewPoW(0, 0)

	block := NewBlock(0, "", "")

	powBlock, err := pow.Mine(block)

	if err != nil {
		t.Fatal(err)
	}

	if powBlock.Get().Nonce != 0 {
		t.Fatal("nonce should not be incremented because PoW difficulty is 0")
	}
}
func TestPoW_Mine_NilBlock(t *testing.T) {
	pow := NewPoW(0, 0)

	_, err := pow.Mine(nil)

	if !errors.Is(err, ErrNilBlock) {
		t.Fatal("should return an nil block error")
	}
}
func TestPoW_Mine_NonceExceeded(t *testing.T) {
	pow := NewPoW(1, POW_LIMIT)

	block := NewBlock(0, "", "")

	block.Get().Nonce = POW_LIMIT

	_, err := pow.Mine(block)

	if !errors.Is(err, ErrNonceExceeded) {
		t.Fatal("should return nonce exceeded error")
	}
}
func TestPoW_Mine(t *testing.T) {
	pow := NewPoW(1, POW_LIMIT)

	block := NewBlock(0, "", "")
	block.Get().Timestamp, _ = time.Parse("2006-01-02T15:04:05.000000Z", "2024-12-07T02:29:21.945508Z")
	block.Get().Hash = NewStringSHA256("")

	powBlock, err := pow.Mine(block)

	if err != nil {
		t.Fatal("should not return an error because mine successfully")
	}
	if powBlock.Get().Nonce != 1 {
		t.Fatal("nonce should be 1", powBlock.Get().Nonce)
	}
}
func TestPoW_Check(t *testing.T) {
	pow := NewPoW(5, 0)

	block := NewBlock(0, "", "")
	block.Get().Hash = "0000f727854b50bb95c054b39c1fe5c92e5ebcfa4bcb5dc279f56aa96a365e5a"

	err := pow.Check(block)

	if !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("should return an invalid proof error, but got %v", err)
	}

	block.Get().Hash = "00000727854b50bb95c054b39c1fe5c92e5ebcfa4bcb5dc279f56aa96a365e5a"

	err = pow.Check(block)

	if err != nil {
		t.Fatalf("should return nil because proof is valid, but got %v", err)
	}
}
