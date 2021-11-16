package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	a         = 1
)

func TCP() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	ID := 0
	for {

		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, ID)
		ID++
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, ID int) {

	for {
		buf := bufio.NewReader(conn)
		data, _ := buf.ReadBytes(10)
		fmt.Println(data, "received from:", conn.RemoteAddr(), "connection number:", ID)
	}
}
