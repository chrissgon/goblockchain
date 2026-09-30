package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// server holds the chain shared by the HTTP handlers. The mutex makes each
// request see and change the chain as one step: /mine reads the last block,
// mines the new one and appends it without another request in between, and
// /chain validates and encodes a chain that no /mine is changing.
type server struct {
	mu    sync.Mutex
	chain IBlockchain
}

func newServer(chain IBlockchain) http.Handler {
	s := &server{chain: chain}

	mux := http.NewServeMux()
	mux.HandleFunc("/mine", s.mine)
	mux.HandleFunc("/chain", s.getChain)

	return mux
}

func (s *server) mine(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.chain.GetChainSize()
	previousBlockHash := s.chain.GetPreviousBlock().Get().Hash
	data := ""

	block := NewBlock(index, previousBlockHash, data)

	_, err := s.chain.Mine(block)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(block)
}

func (s *server) getChain(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.chain.CheckChain()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.chain.GetChain())
}
