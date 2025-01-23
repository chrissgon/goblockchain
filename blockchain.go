package main

import (
	"errors"
)

type Chain []IBlock

type Blockchain struct {
	chain Chain
	miner IMiner
}

type IBlockchain interface {
	GetChain() Chain
	GetMiner() IMiner
	GetPreviousBlock() IBlock
	GetChainSize() int
	Mine(IBlock) (Chain, error)
	CheckChain() error
}

var (
	ErrInvalidChain = errors.New("invalid chain")
	ErrInvalidMiner = errors.New("invalid miner")
)

func NewBlockchain(miner IMiner) (IBlockchain, error) {
	if miner == nil {
		return nil, ErrInvalidMiner
	}

	genesis := NewBlock(0, "", "")

	return &Blockchain{
		chain: []IBlock{genesis},
		miner: miner,
	}, nil
}

func (bc *Blockchain) GetChain() Chain {
	return bc.chain
}
func (bc *Blockchain) GetMiner() IMiner {
	return bc.miner
}
func (bc *Blockchain) Mine(block IBlock) (Chain, error) {
	mineBlock, err := bc.miner.Mine(block)

	if err != nil {
		return bc.chain, err
	}

	err = bc.miner.Check(block)

	if err != nil {
		return bc.chain, err
	}

	err = bc.CheckChain()

	if err != nil {
		return bc.chain, err
	}

	err = block.Check()

	if err != nil {
		return bc.chain, err
	}

	bc.chain = append(bc.chain, mineBlock)

	return bc.chain, nil
}

func (bc *Blockchain) GetPreviousBlock() IBlock {
	return bc.chain[len(bc.chain)-1]
}
func (bc *Blockchain) GetChainSize() int {
	return len(bc.chain)
}
func (bc *Blockchain) CheckChain() error {
	if len(bc.chain) == 1 {
		return nil
	}

	previousBlock := bc.chain[0]

	for i, block := range bc.chain {
		isGenesisBlock := i == 0

		if isGenesisBlock {
			continue
		}

		rightPreviousHash := previousBlock.Get().Hash
		currentPreviousHash := block.Get().PreviousHash

		if currentPreviousHash != rightPreviousHash {
			return ErrInvalidChain
		}

		previousBlock = block
	}

	return nil
}
