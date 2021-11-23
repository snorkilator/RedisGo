package main

import (
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

	toSend := []byte("-" + s.Error() + `\r\n`)
	n, err := conn.Write(toSend) //find out if n can indicate write error (wrong number of bytes printed)
	if err != nil {
		return fmt.Errorf("sendErr: %v", err)
	}
	if n != len(toSend) {
		return fmt.Errorf("sendErr: error message failed to send")
	}
	return nil
}

// Parse takes RESP array of bulk strings, and outputs a slice of []bytes containing the elements of the array
func parse(b []byte) ([][]byte, error) {
	s := string(b) // get rid of this unnecessary copy operation, don't want to have to copy large data if it is unnecessary to do so
	var parsed [][]byte

	if s[0] != byte('*') {
		return nil, fmt.Errorf("parse: input not Array of Bulk Strings type: want *, got %v", string(s[0]))
	}

	arrayLen, currentIndex, err := getLen(1, &s)
	if err != nil {
		return nil, err
	}

	for i := 0; i < arrayLen && currentIndex < len(s); i, currentIndex = i+1, currentIndex+4 {

		if s[currentIndex] != '$' {
			return nil, fmt.Errorf("parse: Expected string symbol '$' for element %v in array but got %v at index %v", i+1, string(s[currentIndex]), currentIndex)
		}
		currentIndex++

		var strLen int
		strLen, currentIndex, err = getLen(currentIndex, &s)
		if err != nil {
			return nil, err
		}

		if currentIndex+strLen > len(s)-1 {
			return nil, fmt.Errorf(`parse: %v indexed element does not exist in array`, i+1)
		}
		tempStr := b[currentIndex : currentIndex+strLen]
		currentIndex += strLen

		if len(s)-1 < currentIndex+3 {
			return nil, fmt.Errorf(`parse: array element %v does not terminate with "\r\n"`, i)
		}
		if s[currentIndex:currentIndex+4] != `\r\n` {
			return nil, fmt.Errorf(`parse: expected \r\n but found %v starting at index %v`, s[currentIndex:currentIndex+4], currentIndex)
		}
		parsed = append(parsed, tempStr)

	}
	return parsed, nil
}

// GetLen is a utility for parse(). It finds the indicator of length for arrays and bulk strings within arrays. Returns length as int as well as the updated index.
func getLen(i int, s *string) (int, int, error) {
	beginStr := i
	var endofNum int

	found := false
	for ; !found; i++ {
		if len(*s)-1 < i+3 {
			return 0, 0, fmt.Errorf(`getLen: Last element in array does not terminate with \r\n`)
		}
		if (*s)[i:i+4] == `\r\n` {
			found = true
			endofNum = i - 1
		}
	}
	lenS := (*s)[beginStr : endofNum+1]

	len, err := strconv.Atoi(lenS)
	if err != nil {
		return 0, 0, fmt.Errorf("getLen: Array length is not valid number: want number but got %v", lenS)
	}

	return len, i + 3, nil
}

// HandleCommand executes command with listed arguments and passes response from command back to caller.
func handleCommand(slc *[][]byte) (resp [][]byte, err error) {
	switch string((*slc)[0]) {
	case "set":
		fmt.Println("set")
		resp, err = set(slc)
	case "get":
		fmt.Println("get")
		resp, err = get(slc)
	default:
		err = fmt.Errorf("handleCommand: no such command %s, try SET or GET", (*slc)[0])
	}
	return
}

// Set puts input into database.
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

// Get gathers and outputs requested data from datbase.
func get(slc *[][]byte) ([][]byte, error) {
	len := len(*slc)
	resp := [][]byte{}
	if len != 2 {
		return nil, fmt.Errorf("get: wrong number of arguments. want 3 but got %v", len)
	}
	v, ok := db[string((*slc)[1])]
	if !ok {
		return nil, fmt.Errorf(`get: "%v" does not exist in database`, (*slc)[1])
	}
	resp = append(resp, v)
	return resp, nil
}

// Respond sends information back to client.
func respond(slc [][]byte, conn net.Conn) error {
	toSend, err := fmtData(slc)
	fmt.Println("sent:", string(toSend))
	conn.Write(toSend)
	return err
}

// FmtData formats input as RESP array of bulk strings.
func fmtData(slc [][]byte) ([]byte, error) {
	delim := `\r\n`
	elCount := fmt.Sprint(len(slc))
	output := []byte(`*` + elCount + delim)

	for _, e := range slc {
		slcLen := fmt.Sprint(len(e))
		elBegin := "$" + slcLen + delim
		output = append(output, []byte(elBegin)...)
		output = append(output, e...)
		output = append(output, []byte(delim)...)
	}
	return output, nil
}
