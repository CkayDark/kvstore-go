package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

func handleConnection(conn net.Conn, store *KVStore) {
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
				conn.Write([]byte("-ERR: wrong number of arguments\n"))
				continue
			}
			store.Set(parts[1], parts[2])
			conn.Write([]byte("+OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("-ERR: parameter missing\n"))
				continue
			}
			val, exists := store.Get(parts[1])
			if exists {
				conn.Write([]byte(val + "\n"))
			} else {
				conn.Write([]byte("(nil)\n"))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("-ERR: parameter missing\n"))
				continue
			}
			ok := store.Delete(parts[1])
			if ok {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}

		default:
			conn.Write([]byte("-ERR: unknown error\n"))
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

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.keyValueMap[key]
	return val, ok
}

func (s *KVStore) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.keyValueMap[key]
	if ok {
		delete(s.keyValueMap, key)
		return ok
	}
	return ok
}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Fehler aufgetreten:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on Port :8080...")
	store := NewKVStore()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Verbindung fehlgeschlagen:", err)
			return
		}

		go handleConnection(conn, store)
	}
}
