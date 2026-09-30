# goblockchain

[Read in English](README.md)

Uma blockchain didática em Go: blocos ligados por hashes SHA-256, um minerador com prova de trabalho e uma API HTTP para minerar blocos e ler a cadeia. Usa só a biblioteca padrão do Go.

O projeto não está terminado. O que ele implementa funciona e tem testes; a lista do que falta está abaixo.

## ⚠️ Requisitos

- Go 1.23.2 ou mais novo (`go.mod`).

## 📦 Instalação

```bash
git clone git@github.com:chrissgon/goblockchain.git
cd goblockchain
```

## 🚀 Início rápido

Rode a API. Ela escuta na porta 8090.

```bash
go run .
```

Em outro terminal, minere um bloco e leia a cadeia:

```bash
curl http://localhost:8090/mine
curl http://localhost:8090/chain
```

`/mine` devolve o bloco novo em JSON, por exemplo:

```json
{"index":1,"nonce":1,"data":"","hash":"0fc4...","previousHash":"14f6...","timestamp":"2026-09-30T03:19:50.096912Z"}
```

`/chain` devolve todos os blocos, começando pelo bloco gênese (índice 0).

## 🧱 O que implementa

| Parte | Onde | O que faz |
|-------|------|-----------|
| Bloco | `block.go`: `NewBlock`, `Block.Check` | Um bloco tem índice, nonce, dados, hash, hash anterior e um horário em UTC. `Check` recalcula o hash e devolve `invalid block hash` se ele não bater. |
| Hash | `helpers.go`: `GetBlockData`, `NewStringSHA256` | O hash é o SHA-256 do JSON do bloco sem o campo de hash. |
| Cadeia | `blockchain.go`: `NewBlockchain`, `Mine`, `CheckChain` | Começa com um bloco gênese. `Mine` roda o minerador, confere a prova, as ligações da cadeia e o hash do bloco, e então acrescenta o bloco. `CheckChain` confere se o hash anterior de cada bloco é igual ao hash do bloco antes dele. |
| Prova de trabalho | `miner.go`: `NewPoW`, `PoW.Mine`, `PoW.Check` | Incrementa o nonce até o hash começar com `difficulty` zeros, e para com `nonce limit exceeded` quando o nonce chega ao limite. |
| API HTTP | `main.go`, `server.go`: `newServer` | `GET /mine` minera um bloco com dados vazios; `GET /chain` valida e devolve a cadeia. O minerador é criado com dificuldade 1 e limite de nonce 50. Um mutex em `server` atende um pedido de cada vez, então chamadas simultâneas a `/mine` acrescentam cada uma um bloco ligado ao anterior. |

## 🚧 O que falta

- **Transações.** `/mine` sempre cria um bloco com dados vazios (`data := ""` em `server.go`); não há como enviar dados ou transações.
- **Rede entre nós.** Existe uma única cadeia em memória; não há pares, troca de cadeia nem consenso entre nós.
- **Persistência.** A cadeia fica em memória e se perde quando o processo para.
- **Validação completa da cadeia.** `CheckChain` confere só as ligações de hash anterior; não recalcula o hash de cada bloco nem a prova de trabalho.
- **A mineração pode falhar.** Com dificuldade 1 e limite de nonce 50, um bloco tem cerca de 4% de chance ((15/16)^50) de não achar um hash válido; `/mine` então responde 500 com `nonce limit exceeded`.
- **Configuração.** Dificuldade, limite de nonce e porta são fixos em `main.go`.
- `ErrEmptyBlockData` está declarado em `miner.go`, mas nunca é usado.

## 📝 Testes

Rode os testes, com o detector de corrida:

```bash
go test -race ./...
```

`TestServer_Concurrent` em `server_test.go` envia `/mine` e `/chain` ao mesmo tempo.

Rode os testes com cobertura e abra o relatório em HTML:

```bash
go test ./... -coverprofile=coverage.out && ./coverage-ignore.sh && go tool cover -html=coverage.out
```

## 💪🏻 Contribuição

Este projeto é de código aberto e aceita contribuições da comunidade. Fique à vontade para fazer um fork, implementar melhorias e enviar um pull request.

## ❤️ Autores

- [@chrissgon](https://www.github.com/chrissgon)
