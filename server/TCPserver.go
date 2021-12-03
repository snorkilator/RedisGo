package server

import (
	"log"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type server struct {
	l net.Listener
}

var Server = server{*new(net.Listener)}

func (Server server) Close() {
	Server.l.Close()
}

// Run receives tcp connections and passes them to handler in seperate thread
func Run(msgCh chan ClientMHandle) {
	// Listen for incoming connections.

	var err error
	Server.l, err = net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer Server.l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	ID := 0
	for {

		// Listen for an incoming connection.
		conn, err := Server.l.Accept()
		if err != nil {
			//log.Println("Error accepting: ", err.Error())
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
		buf := make([]byte, 1000)
		_, err := conn.Read(buf)
		if err != nil {
			log.Println("handleRequest:", err)
			conn.Close()
			return
		}
		clientMHandle := ClientMHandle{buf, conn}
		msgCh <- clientMHandle
		log.Printf("handleRequest: %v received from: %v connection number: %v", string(buf), conn.RemoteAddr(), ID)
	}
}
