package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	message := "Willkommen beim KV-Store!\n"
	conn.Write([]byte(message))
	for scanner.Scan() {
		message := scanner.Text()
		conn.Write([]byte("You said: " + message + "\n"))
	}
}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Fehler aufgetreten:", err)
		return
	}
	defer listener.Close()

	for {
		fmt.Println("Server is listening on Port :8080...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Verbindung fehlgeschlagen:", err)
			return
		}

		handleConnection(conn)
	}
}
