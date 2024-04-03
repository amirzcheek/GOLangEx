package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	CONN_PORT = ":3335"
	CONN_TYPE = "tcp"
)

var wg sync.WaitGroup

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(message)
	}
	wg.Done()
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)
	fmt.Print("Enter your name please: ")
	name, _ := reader.ReadString('\n')
	_, err := writer.WriteString(name)
	if err != nil {
		wg.Done()
	}

	err = writer.Flush()
	if err != nil {
		wg.Done()
	}

	for {
		message, _ := reader.ReadString('\n')
		_, err := writer.WriteString(message)
		if err != nil || strings.TrimSpace(message) == "/quit" {
			break
		}
		err = writer.Flush()
		if err != nil {
			break
		}
	}
	wg.Done()
}

func main() {
	wg.Add(2)

	conn, err := net.Dial(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go Read(conn)
	go Write(conn)

	wg.Wait()
}
