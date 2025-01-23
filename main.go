package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	// os.Setenv("POW_LIMIT", "50")

	miner := NewPoW(1, 50)
	chain, _ := NewBlockchain(miner)

	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {

		index := chain.GetChainSize()
		previousBlockHash := chain.GetPreviousBlock().Get().Hash
		data := ""

		block := NewBlock(index, previousBlockHash, data)

		_, err := chain.Mine(block)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(block)
	})

	http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		err := chain.CheckChain()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chain.GetChain())
	})

	http.ListenAndServe(":8090", nil)
}
