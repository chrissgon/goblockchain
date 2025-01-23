package main

import (
	"errors"
	"time"
)

type Block struct {
	Index        int       `json:"index"`
	Nonce        int       `json:"nonce"`
	Data         any       `json:"data"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previousHash"`
	Timestamp    time.Time `json:"timestamp"`
}

type IBlock interface {
	Get() *Block
	Check() error
}

var (
	ErrInvalidBlockHash = errors.New("invalid block hash")
)

func NewBlock(index int, previousHash string, data any) IBlock {
	block := Block{
		Index:        index,
		Nonce:        0,
		Data:         data,
		PreviousHash: previousHash,
		Timestamp:    time.Now().UTC(),
	}

	block.Hash = NewStringSHA256(GetBlockData(block))

	return &block
}

func (b *Block) Get() *Block {
	return b
}
func (b *Block) Check() error {
	rightHash := NewStringSHA256(GetBlockData(*b))
	currentHash := b.Hash

	if currentHash != rightHash {
		return ErrInvalidBlockHash
	}

	return nil
}
