package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestServer_Mine(t *testing.T) {
	chain, _ := NewBlockchain(NewPoW(0, 0))
	srv := httptest.NewServer(newServer(chain))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/mine")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status should be 200, have %d", res.StatusCode)
	}

	var block Block
	if err := json.NewDecoder(res.Body).Decode(&block); err != nil {
		t.Fatal(err)
	}
	if block.Index != 1 {
		t.Fatalf("mined block index should be 1, have %d", block.Index)
	}
}

func TestServer_MineError(t *testing.T) {
	chain, _ := NewBlockchain(NewPoW(1, 0))
	srv := httptest.NewServer(newServer(chain))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/mine")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status should be 500, have %d", res.StatusCode)
	}
}

func TestServer_ChainError(t *testing.T) {
	chain, _ := NewBlockchain(NewPoW(0, 0))
	bc := chain.(*Blockchain)
	bc.chain = append(bc.chain, NewBlock(1, "not-the-genesis-hash", ""))

	srv := httptest.NewServer(newServer(chain))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/chain")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status should be 500, have %d", res.StatusCode)
	}
}

// TestServer_Concurrent sends /mine and /chain at the same time. Run it with
// `go test -race ./...`: without the lock in server the race detector reports
// the shared chain, and blocks mined at the same time link to the same
// previous block, which breaks the chain.
func TestServer_Concurrent(t *testing.T) {
	const requests = 20

	chain, _ := NewBlockchain(NewPoW(0, 0))
	srv := httptest.NewServer(newServer(chain))
	defer srv.Close()

	var wg sync.WaitGroup
	errs := make(chan error, 2*requests)

	for range requests {
		for _, path := range []string{"/mine", "/chain"} {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				res, err := http.Get(srv.URL + path)
				if err != nil {
					errs <- err
					return
				}
				res.Body.Close()
				if res.StatusCode != http.StatusOK {
					errs <- &statusError{path: path, code: res.StatusCode}
				}
			}(path)
		}
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}

	if chain.GetChainSize() != requests+1 {
		t.Fatalf("chain size should be %d, have %d", requests+1, chain.GetChainSize())
	}
	if err := chain.CheckChain(); err != nil {
		t.Fatalf("chain should be valid: %v", err)
	}
}

type statusError struct {
	path string
	code int
}

func (e *statusError) Error() string {
	return e.path + " answered " + http.StatusText(e.code)
}
