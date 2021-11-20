package server

import (
	"bufio"
	"log"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

// Run receives tcp connections and passes them to handler in seperate thread
func Run(msgCh chan ClientMHandle) {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	ID := 0
	for {

		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, ID, msgCh)
		ID++
	}
}

type ClientMHandle = struct {
	Data []byte
	Conn net.Conn
}

// handleRequests Handles incoming requests and sends incoming data to chan for processing.
func handleRequest(conn net.Conn, ID int, msgCh chan ClientMHandle) {

	for {
		buf := bufio.NewReader(conn)
		data, err := buf.ReadBytes(10)
		if err != nil {
			log.Println("handleRequest:", err)
			return
		}
		clientMHandle := ClientMHandle{data, conn}
		msgCh <- clientMHandle
		log.Println("handleRequest:", string(data), "received from:", conn.RemoteAddr(), "connection number:", ID)
	}
}
