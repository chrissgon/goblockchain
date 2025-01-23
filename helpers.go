package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

func NewStringSHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	bin := hash.Sum(nil)
	return fmt.Sprintf("%x", bin)
}

func GetBlockData(block Block) string {
	type BlockData struct {
		Index        int       `json:"index"`
		Nonce        int       `json:"nonce"`
		Data         any       `json:"data"`
		PreviousHash string    `json:"previousHash"`
		Timestamp    time.Time `json:"timestamp"`
	}

	treatedBlock := BlockData{
		Index:        block.Index,
		Nonce:        block.Nonce,
		Data:         block.Data,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
	}

	blockBytes, _ := json.Marshal(treatedBlock)

	return string(blockBytes)
}
