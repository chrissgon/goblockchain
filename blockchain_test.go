package main

import (
	"errors"
	"testing"
)

func TestBlockchain_New(t *testing.T) {
	_, err := NewBlockchain(nil)

	if !errors.Is(err, ErrInvalidMiner) {
		t.Fatal("should return an error because miner is nil")
	}

	pow := NewPoW(0, 0)

	chain, err := NewBlockchain(pow)

	if err != nil {
		t.Fatal("chain should be create successfully")
	}

	if len(chain.GetChain()) != 1 {
		t.Fatal("chain should start with genesis block")
	}

	if chain.GetMiner() == nil {
		t.Fatal("chain should start with a miner set")
	}
}
func TestBlockchain_GetChain(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	if chain.GetChain() == nil {
		t.Fatal("chain should not be nil")
	}
}
func TestBlockchain_GetMiner(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	if chain.GetMiner() == nil {
		t.Fatal("miner should not be nil")
	}
}
func TestBlockchain_GetChainSize(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	if chain.GetChainSize() != 1 {
		t.Fatal("chain size should be 1")
	}
}
func TestBlockchain_GetPreviousBlock(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	if chain.GetPreviousBlock().Get().Index != 0 {
		t.Fatal("should return the right previous block")
	}
}

func TestBlockchain_Mine_Miner(t *testing.T) {
	pow := NewPoW(1, 0)

	chain, _ := NewBlockchain(pow)

	index := chain.GetChainSize()
	previousBlockHash := chain.GetPreviousBlock().Get().Hash

	block := NewBlock(index, previousBlockHash, "")

	_, err := chain.Mine(block)

	if !errors.Is(err, ErrNonceExceeded) {
		t.Fatal("should return a nonce exceeded error")
	}
}
func TestBlockchain_Mine_CheckChain(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	block := NewBlock(0, "", "")

	chain.Mine(block)
	_, err := chain.Mine(block)

	if !errors.Is(err, ErrInvalidChain) {
		t.Fatal("should return a invalid chain error")
	}
}
func TestBlockchain_Mine(t *testing.T) {
	pow := NewPoW(0, 0)

	chain, _ := NewBlockchain(pow)

	for range 3 {
		index := chain.GetChainSize()
		previousBlockHash := chain.GetPreviousBlock().Get().Hash

		block := NewBlock(index, previousBlockHash, "")

		_, err := chain.Mine(block)

		if err != nil {
			t.Fatal("should not return an error because mine successfully")
		}

	}

	index := chain.GetChainSize()
	previousBlockHash := chain.GetChain()[2].Get().Hash

	if chain.GetChainSize() != 4 {
		t.Fatal("chain size should be 4")
	}
	if chain.GetPreviousBlock().Get().Index != index-1 {
		t.Fatal("previous block index is invalid")
	}
	if chain.GetPreviousBlock().Get().PreviousHash != previousBlockHash {
		t.Fatal("previous block hash is invalid")
	}
}
