# goblockchain

A small educational blockchain in Go: blocks linked by SHA-256 hashes, a proof-of-work miner, and an HTTP API to mine blocks and read the chain. It uses only the Go standard library.

The project is not finished. What it implements works and is covered by tests; the list of what is missing is below.

## ⚠️ Requirements

- Go 1.23.2 or later (`go.mod`).

## 📦 Install

```bash
git clone git@github.com:chrissgon/goblockchain.git
cd goblockchain
```

## 🚀 Quick Start

Run the API. It listens on port 8090.

```bash
go run .
```

In another terminal, mine a block and read the chain:

```bash
curl http://localhost:8090/mine
curl http://localhost:8090/chain
```

`/mine` returns the new block as JSON, for example:

```json
{"index":1,"nonce":1,"data":"","hash":"0fc4...","previousHash":"14f6...","timestamp":"2026-09-30T03:19:50.096912Z"}
```

`/chain` returns every block, starting with the genesis block (index 0).

## 🧱 What it implements

| Piece | Where | What it does |
|-------|-------|--------------|
| Block | `block.go`: `NewBlock`, `Block.Check` | A block has index, nonce, data, hash, previous hash and a UTC timestamp. `Check` recomputes the hash and returns `invalid block hash` if it does not match. |
| Hashing | `helpers.go`: `GetBlockData`, `NewStringSHA256` | The hash is the SHA-256 of the block's JSON without the hash field. |
| Chain | `blockchain.go`: `NewBlockchain`, `Mine`, `CheckChain` | Starts with a genesis block. `Mine` runs the miner, checks the proof, the chain links and the block hash, then appends the block. `CheckChain` verifies that each block's previous hash equals the hash of the block before it. |
| Proof of work | `miner.go`: `NewPoW`, `PoW.Mine`, `PoW.Check` | Increments the nonce until the hash starts with `difficulty` zeros, and stops with `nonce limit exceeded` when the nonce reaches the limit. |
| HTTP API | `main.go`, `server.go`: `newServer` | `GET /mine` mines a block with empty data; `GET /chain` validates and returns the chain. The miner is created with difficulty 1 and a nonce limit of 50. A mutex in `server` serializes the requests, so concurrent `/mine` calls each append one block linked to the one before. |

## 🚧 What is missing

- **Transactions.** `/mine` always creates a block with empty data (`data := ""` in `main.go`); there is no way to submit data or transactions.
- **Network between nodes.** There is a single in-memory chain; no peers, no chain exchange, no consensus between nodes.
- **Persistence.** The chain lives in memory and is lost when the process stops.
- **Full chain validation.** `CheckChain` checks only the previous-hash links; it does not recompute each block's hash or its proof of work.
- **Mining can fail.** With difficulty 1 and a nonce limit of 50, a block has about a 4% chance ((15/16)^50) of not finding a valid hash; `/mine` then answers 500 with `nonce limit exceeded`.
- **Configuration.** Difficulty, nonce limit and port are fixed in `main.go`.
- `ErrEmptyBlockData` is declared in `miner.go` but never used.

## 📝 Tests

Run the tests, with the race detector:

```bash
go test -race ./...
```

`TestServer_Concurrent` in `server_test.go` sends `/mine` and `/chain` at the same time.

Run the tests with coverage and open the HTML report:

```bash
go test ./... -coverprofile=coverage.out && ./coverage-ignore.sh && go tool cover -html=coverage.out
```

## 💪🏻 Contribution

This project is open source and welcomes community contributions. Feel free to fork, implement improvements, and submit a pull request.

## ❤️ Authors

- [@chrissgon](https://www.github.com/chrissgon)
