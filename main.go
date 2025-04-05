// main.go
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	// global in‑memory store, protected by a mutex
	store = make(map[string]string)
	mu    sync.RWMutex
)

func main() {
	log.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected: %s", conn.RemoteAddr())

	resp := NewResp(conn)
	writer := NewWriter(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			log.Printf("read error from %s: %v", conn.RemoteAddr(), err)
			return
		}

		// Expect an ARRAY of BULKs for every Redis command
		if value.typ != "array" {
			writer.Write(Value{typ: "error", str: "ERR protocol error: expected array"})
			continue
		}

		// Log the raw parsed command
		log.Printf("command: %#v", value.array)

		// At least one element (the command name) is required
		if len(value.array) < 1 || value.array[0].typ != "bulk" {
			writer.Write(Value{typ: "error", str: "ERR protocol error: bad command format"})
			continue
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		switch cmd {
		case "PING":
			writer.Write(Value{typ: "string", str: "PONG"})

		case "SET":
			if len(args) < 2 || args[0].typ != "bulk" || args[1].typ != "bulk" {
				writer.Write(Value{typ: "error", str: "ERR wrong number of arguments for 'SET'"})
				continue
			}
			key := args[0].bulk
			val := args[1].bulk

			mu.Lock()
			store[key] = val
			mu.Unlock()

			writer.Write(Value{typ: "string", str: "OK"})

		case "GET":
			if len(args) < 1 || args[0].typ != "bulk" {
				writer.Write(Value{typ: "error", str: "ERR wrong number of arguments for 'GET'"})
				continue
			}
			key := args[0].bulk

			mu.RLock()
			val, ok := store[key]
			mu.RUnlock()

			if !ok {
				writer.Write(Value{typ: "null"})
			} else {
				writer.Write(Value{typ: "bulk", bulk: val})
			}

		default:
			writer.Write(Value{typ: "error", str: fmt.Sprintf("ERR unknown command '%s'", cmd)})
		}
	}
}
