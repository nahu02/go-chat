# Go WebSocket Chatroom

Simple, fully-functional chat server written in **Go 1.22**.
It showcases how to:

* expose HTTP endpoints and upgrade them to WebSocket connections
* broadcast messages to dynamic “rooms” with minimal state
* handle concurrent clients safely in pure Go without external message brokers

The accompanying frontend (HTML + CSS) was generated with AI to keep the focus on backend learning.

## Why I built it

1. **Learn Go’s HTTP stack** – handlers, simple routing, file serving.
2. **Explore concurrency primitives** – goroutines and sharing state across connections.
3. **Work with WebSockets** – real-time, bi-directional communication.
4. **Produce a deployable artefact** – something that can be cloned, ran and easily seen in action.

I deliberately **avoided AI code generation** or copy-pasting; every Go line is hand-written after reading the relevant Gorilla WebSocket docs and one tutorial article for initial direction.

## How it works

| Component | Responsibility |
| :-- | :-- |
| `main.go` | Entry point. Registers HTTP routes, starts server on port 8080. |
| `roomHandler` | Upgrades any request to a WebSocket connection, then hands it to the room handler. |
| `handleRoomConnection` | Runs in its own goroutine per client. Reads JSON messages, broadcasts them to every connection stored in `roomClients[roomId]`. |
| `roomClients` | In-memory map `map[string][]*websocket.Conn` keeping active connections per room. |
| Frontend (`public/`) | Minimal page that opens a WebSocket, lets users join a room and exchange messages. |

Concurrency safety is achieved by:

* creating a new goroutine per connection
* a global `sync.Mutex` (`roomClientsMutex`) guards all reads / writes to the `roomClients` map
* cleaning up closed connections with the helper `removeFromRoom`

For a production-grade system I would store the rooms and messages in a specialized external service (mainly for fault tolerance, traceability, and to fine-grain the concurrency locking on a per-room basis), but for a learning project the current setup is enough to demonstrate concepts.

## Quick start

```bash
# 1. clone the repo
git clone https://github.com/nahu02/go-chat
cd go-chat

# 2. download dependencies
go mod download

# 3. run
go run .

# 4. open several tabs
http://localhost:8080
# join the same room ID in each tab and chat in real-time
```

_No database or message queue needed. Terminate with `Ctrl+C`._

## What I learned

* Building an HTTP server and static file server in Go in <30 lines.
* Upgrading HTTP to WebSocket with Gorilla and handling JSON payloads.
* Managing a simple shared state across goroutines and cleaning up resources.
* Structured logging for both incoming messages and error paths.
* Writing readable, idiomatic Go: error handling, slices, maps and naming.


## Possible next steps

* Use an external data storage for the chat rooms.
* Switch from the single `sync.Mutex` to either a per-room mutex.
* Consider using `sync.Map` instead of the `sync.Mutex`.
* Implement CORS / origin checks in the upgrader.
* Introduce TLS for secure WebSocket (`wss://`) communication.
* Containerize with a small `Dockerfile` for easy cloud deployment.
* Write unit tests for the room logic.


### A note on inspiration

The overall approach was inspired by a [tutorial on implementing WebSockets in Go](https://medium.com/wisemonks/implementing-websockets-in-golang-d3e8e219733b), but the final implementation, refactoring and comments are entirely my own interpretation and exercise.

