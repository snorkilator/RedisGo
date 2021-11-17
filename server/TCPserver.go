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
)

// receives tcp connections and passes them to handler in seperate thread
func Run(a chan ClientMHandle) {
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
		go handleRequest(conn, ID, a)
		ID++
	}
}

type ClientMHandle = struct {
	Data []byte
	Conn net.Conn
}

// Handles incoming requests.
func handleRequest(conn net.Conn, ID int, a chan ClientMHandle) {

	for {
		buf := bufio.NewReader(conn)
		data, _ := buf.ReadBytes(10)
		ClientMHandle := ClientMHandle{data, conn}
		a <- ClientMHandle
		fmt.Println("server.Run():", string(data), "received from:", conn.RemoteAddr(), "connection number:", ID)
	}
}
