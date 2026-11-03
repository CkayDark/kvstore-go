package main

import (
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Fehler aufgetreten:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on Port :8080...")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Verbindung fehlgeschlagen:", err)
		return
	}
	defer conn.Close()

	message := "Willkommen beim KV-Store!\n"
	conn.Write([]byte(message))
}
