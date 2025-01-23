package main

import (
	"errors"
	"strings"
)

var (
	ErrNonceExceeded  = errors.New("nonce limit exceeded")
	ErrInvalidProof   = errors.New("invalid proof")
	ErrNilBlock       = errors.New("nil block")
	ErrEmptyBlockData = errors.New("empty block data")
)

type IMiner interface {
	Mine(IBlock) (IBlock, error)
	Check(IBlock) error
}

type PoW struct {
	difficulty int
	limit      int
}

func NewPoW(difficulty int, limit int) IMiner {
	return &PoW{difficulty: difficulty, limit: limit}
}

func (pow *PoW) Mine(block IBlock) (IBlock, error) {
	if block == nil {
		return block, ErrNilBlock
	}

	data := GetBlockData(*block.Get())

	block.Get().Hash = NewStringSHA256(data)

	if pow.difficulty == 0 {
		return block, nil
	}

	nonce := block.Get().Nonce

	if nonce >= pow.limit {
		return block, ErrNonceExceeded
	}

	err := pow.Check(block)

	if err != nil {
		nonce++
		block.Get().Nonce = nonce

		return pow.Mine(block)
	}

	return block, nil
}

func (pow *PoW) Check(block IBlock) error {
	proof := block.Get().Hash[:pow.difficulty] == strings.Repeat("0", pow.difficulty)

	if proof {
		return nil
	}

	return ErrInvalidProof
}
