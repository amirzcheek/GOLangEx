package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	CONN_PORT = ":3335"
	CONN_TYPE = "tcp"
)

var (
	clients    = make(map[string]net.Conn)
	clientsMux sync.Mutex
)

func broadcast(message string, sender string) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for name, conn := range clients {
		if name != sender {
			fmt.Fprintf(conn, "[%s]: %s\n", sender, message)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientsMux.Lock()
	clients[name] = conn
	clientsMux.Unlock()

	fmt.Printf("[%s] connected\n", name)
	broadcast(fmt.Sprintf("[%s] joined the chat", name), name)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message == "/quit" {
			break
		}
		broadcast(message, name)
	}

	clientsMux.Lock()
	delete(clients, name)
	clientsMux.Unlock()

	fmt.Printf("[%s] disconnected\n", name)
	broadcast(fmt.Sprintf("[%s] left the chat", name), name)
}

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on " + CONN_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		go handleClient(conn)
	}
}
