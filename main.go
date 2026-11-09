package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	conn.Write([]byte("Willkommen beim KV-Store!\n"))

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		parts[0] = strings.ToUpper(parts[0])

		switch parts[0] {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERR: Parameter fehlen\n"))
				continue
			}
			conn.Write([]byte("Benutzer erfolgreich erstellt\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERR: Parameter fehlt\n"))
				continue
			}
			conn.Write([]byte("Eintrag\n"))

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERR: Parameter fehlt\n"))
				continue
			}
			conn.Write([]byte("Eintrag gelöscht\n"))

		default:
			conn.Write([]byte("ERR: Unbekannter Befehl\n"))
		}
	}
}

type KVStore struct {
	keyValueMap map[string]string
	mu          sync.RWMutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		keyValueMap: make(map[string]string),
		mu:          sync.RWMutex{},
	}
}

func (s *KVStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyValueMap[key] = value
}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Fehler aufgetreten:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on Port :8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Verbindung fehlgeschlagen:", err)
			return
		}

		go handleConnection(conn)
	}
}
