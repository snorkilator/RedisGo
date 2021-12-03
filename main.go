package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"redis/server"
	"strconv"
)

var db map[string][]byte

// setCommand = *3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n
// getCommand = *2\r\n$3\r\nget\r\n$1\r\na\r\n

func main() {
	messageCh := make(chan server.ClientMHandle)
	db = make(map[string][]byte)

	go server.Run(messageCh)
	for msg := range messageCh {

		s, err := parse(msg.Data)
		if err != nil {
			log.Println(err)
			err = sendErr(err, msg.Conn)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		resp, err := handleCommand(&s)
		if err != nil {
			log.Println(err)
			err = sendErr(err, msg.Conn)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		err = respond(resp, msg.Conn)
		if err != nil {
			log.Println(err)
			err = sendErr(err, msg.Conn)
			if err != nil {
				log.Println(err)
			}
			continue
		}
	}
}

func sendErr(s error, conn net.Conn) error {

	toSend := []byte("-" + s.Error())
	toSend = append(toSend, []byte("\r\n")...)
	n, err := conn.Write(toSend) //find out if n can indicate write error (wrong number of bytes printed)
	if err != nil {
		return fmt.Errorf("sendErr: %v", err)
	}
	if n != len(toSend) {
		return fmt.Errorf("sendErr: error message failed to send")
	}
	return nil
}

// parse takes RESP array of bulk strings, and outputs a slice of []bytes containing the elements of the array
func parse(b []byte) ([][]byte, error) {
	var parsed [][]byte
	if b[0] != '*' {
		return nil, fmt.Errorf("parse: input not Array of Bulk Strings type: want *, got %v", string(b[0]))
	}

	arrayLen, currentIndex, err := getLen(1, &b)
	if err != nil {
		return nil, err
	}

	for i := 0; i < arrayLen && currentIndex < len(b); i, currentIndex = i+1, currentIndex+2 {

		if b[currentIndex] != '$' {
			return nil, fmt.Errorf("parse: Expected string symbol '$' for element %v in array but got %v at index %v", i+1, string(b[currentIndex]), currentIndex)
		}
		currentIndex++

		var strLen int
		strLen, currentIndex, err = getLen(currentIndex, &b)
		if err != nil {
			return nil, err
		}

		if currentIndex+strLen > len(b)-1 {
			return nil, fmt.Errorf(`parse: %v indexed element does not exist in array`, i+1)
		}
		tempStr := b[currentIndex : currentIndex+strLen]
		currentIndex += strLen

		if len(b)-1 < currentIndex+1 {
			return nil, fmt.Errorf(`parse: array element %v does not terminate with "\r\n"`, i)
		}
		if !bytes.Equal(b[currentIndex:currentIndex+2], []byte("\r\n")) {
			return nil, fmt.Errorf(`parse: expected \r\n but found %v starting at index %v`, b[currentIndex:currentIndex+4], currentIndex)
		}
		parsed = append(parsed, tempStr)

	}
	return parsed, nil
}

// getLen is a utility for parse(). It finds the indicator of length for arrays and bulk strings within arrays. Returns length as int as well as the updated index.
func getLen(i int, b *[]byte) (int, int, error) {
	beginStr := i
	var endofNum int

	found := false
	for ; !found; i++ {
		if len(*b) < i+2 {
			return 0, 0, fmt.Errorf(`getLen: Last element in array does not terminate with \r\n`)
		}
		if bytes.Equal((*b)[i:i+2], []byte("\r\n")) {
			found = true
			endofNum = i - 1
		}
	}
	lenB := (*b)[beginStr : endofNum+1]

	len, err := strconv.Atoi(string(lenB))
	if err != nil {
		return 0, 0, fmt.Errorf("getLen: Array length is not valid number: want number but got %v", lenB)
	}

	return len, i + 1, nil
}

// handleCommand executes command with listed arguments and passes response from command back to caller.
func handleCommand(slc *[][]byte) (resp [][]byte, err error) {
	switch string((*slc)[0]) {
	case "set":
		fmt.Println("set")
		resp, err = set(slc)
	case "get":
		fmt.Println("get")
		resp, err = get(slc)
	case "info":
		resp = [][]byte{[]byte("OK")}
	default:
		err = fmt.Errorf("handleCommand: no such command %s, try set or get", (*slc)[0])
	}
	return
}

// set puts input into database.
func set(slc *[][]byte) ([][]byte, error) {
	resp := [][]byte{}
	len := len(*slc)
	if len != 3 {
		return nil, fmt.Errorf("set: wrong number of arguments. want 3 but got %v", len)
	}
	db[string((*slc)[1])] = (*slc)[2]

	resp = append(resp, []byte("OK"))
	return resp, nil
}

// get collects and outputs requested data.
func get(slc *[][]byte) ([][]byte, error) {
	len := len(*slc)
	resp := [][]byte{}
	if len != 2 {
		return nil, fmt.Errorf("get: wrong number of arguments. want 2 but got %v", len)
	}
	v, ok := db[string((*slc)[1])]
	if !ok {
		return nil, fmt.Errorf(`get: "%v" does not exist in database`, (*slc)[1])
	}
	resp = append(resp, v)
	return resp, nil
}

// respond sends information back to client.
func respond(slc [][]byte, conn net.Conn) error {
	toSend := fmtData(slc)
	fmt.Println("sent:", string(toSend))
	n, err := conn.Write(toSend)
	if err != nil {
		return fmt.Errorf("sendErr: %v", err)
	}
	if n != len(toSend) {
		return fmt.Errorf("sendErr: error message failed to send")
	}
	return err
}

// fmtData formats input as RESP array of bulk strings.
func fmtData(slc [][]byte) []byte {
	delim := []byte("\r\n")
	elCount := fmt.Sprint(len(slc))
	var output []byte
	output = append(output, '*')
	output = append(output, []byte(elCount)...)
	output = append(output, delim...)

	for _, e := range slc {
		slcLen := fmt.Sprint(len(e))
		var elBegin = []byte("$")
		elBegin = append(elBegin, slcLen...)
		elBegin = append(elBegin, delim...)
		output = append(output, []byte(elBegin)...)
		output = append(output, e...)
		output = append(output, delim...)
	}
	return output
}
