# GO_Redis_Server
Build Your Own Redis Server from Scratch 

# Go Redis Clone

A minimal, in‑memory Redis‑like server written in Go.  
Supports the RESP protocol with basic commands: **PING**, **SET**, and **GET**.  
Ideal as a learning project for TCP servers, RESP parsing, and Go concurrency.

---

## 🚀 Features

- **RESP parser & writer**  
- **Commands**:  
  - `PING` → `PONG`  
  - `SET <key> <value>` → `OK`  
  - `GET <key>` → bulk string or `(nil)`  
- **Concurrent clients**: one goroutine per connection  
- **Thread‑safe** in‑memory store (`map[string]string` + `sync.RWMutex`)  
- **Logging** of connections, commands, and errors  

---

## 🛠️ Prerequisites

- [Go 1.16+](https://golang.org/dl/)  
- A Redis client for testing (e.g. `redis-cli`)

---

## 📦 Installation & Running

```bash
# Clone this repo
git clone https://github.com/yourusername/go-redis-clone.git
cd go-redis-clone

# (Optional) Initialize Go module
go mod init github.com/yourusername/go-redis-clone

# Build the server binary
go build -o redis-server main.go resp.go

# Run the server
./redis-server

# Or run directly with go run
go run main.go resp.go
------------------------------------------------------------
 💬 Usage
Open a new terminal and connect with redis-cli:


$ redis-cli
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET foo bar
OK
127.0.0.1:6379> GET foo
"bar"
127.0.0.1:6379> GET missing
(nil)
127.0.0.1:6379> UNKNOWN
(error) ERR unknown command 'UNKNOWN'
