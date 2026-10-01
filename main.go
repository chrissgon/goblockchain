package main

import (
	"net/http"
)

func main() {
	// os.Setenv("POW_LIMIT", "50")

	miner := NewPoW(1, 50)
	chain, _ := NewBlockchain(miner)

	http.ListenAndServe(":8090", newServer(chain))
}
