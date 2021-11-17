package main

import (
	"fmt"
	"net"
	"redis/server"
	"strconv"
)

var db map[string][]byte

// setCommand = *3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n
// getCommand = *2\r\n$3\r\nget\r\n$1\r\na\r\n

func main() {
	a := make(chan server.ClientMHandle)
	db = make(map[string][]byte)

	go server.Run(a)
	for {
		cMessage := <-a

		s, err := parser(cMessage.Data)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := commandHandler(&s)
		if err != nil {
			// add error sender
			fmt.Println(err)
		}

		err = responder(resp, cMessage.Conn)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//Takes RESP array of bulk strings, and outputs a slice of strings containing the elements of the array
func parser(b []byte) ([][]byte, error) {
	s := string(b) // get rid of this unnecessary copy operation, don't want to have to copy large data if it is unnecessary to do so
	var parsed [][]byte

	if s[0] != byte('*') {
		return nil, fmt.Errorf("input is not Array of Bulk Strings type: want *, got %v", string(s[0]))
	}

	arrayLen, currentIndex, err := getLen(1, &s)
	if err != nil {
		return nil, err
	}

	for i := 0; i < arrayLen && currentIndex < len(s); i, currentIndex = i+1, currentIndex+4 {

		if s[currentIndex] != '$' {
			return nil, fmt.Errorf("parser: Expected type symbol '$' for element %v in array but got %v at index %v", i+1, string(s[currentIndex]), currentIndex)
		}
		currentIndex++

		var strLen int
		strLen, currentIndex, err = getLen(currentIndex, &s)
		if err != nil {
			return nil, err
		}

		if currentIndex+strLen > len(s)-1 {
			return nil, fmt.Errorf(`parser: %v indexed element does not exist in array`, i+1)
		}
		tempStr := b[currentIndex : currentIndex+strLen]
		currentIndex += strLen

		if len(s)-1 < currentIndex+3 {
			return nil, fmt.Errorf(`parser: array element %v does not terminate with \r\n`, i)
		}
		if s[currentIndex:currentIndex+4] != `\r\n` {
			return nil, fmt.Errorf(`parser: expected \r\n but found %v starting at index %v`, s[currentIndex:currentIndex+4], currentIndex)
		}
		parsed = append(parsed, tempStr)

	}
	return parsed, nil
}

//Starts at some index and checks for \r\n. If it finds that, it will return a string starting at the endex and ending before the \r\n.
//Also returns index value after \r\n
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

//reads first element of slc and executes command found there, or returns error
// passes response from command back to caller
func commandHandler(slc *[][]byte) (resp [][]byte, err error) {
	switch string((*slc)[0]) {
	case "set":
		fmt.Println("set")
		resp, err = set(slc)
	case "get":
		fmt.Println("get")
		resp, err = get(slc)
	default:
		err = fmt.Errorf("No such command " + string((*slc)[0]) + ", try SET or GET")
	}
	return
}

// puts input into db
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

// finds and sends input from db
func get(slc *[][]byte) ([][]byte, error) {
	len := len(*slc)
	resp := [][]byte{}
	if len != 2 {
		return nil, fmt.Errorf("set: wrong number of arguments. want 3 but got %v", len)
	}
	v, ok := db[string((*slc)[1])]
	if !ok {
		return nil, fmt.Errorf("%v does not exist in database", (*slc)[1])
	}
	resp = append(resp, v)
	return resp, nil
}

//sends information back to client
func responder(slc [][]byte, conn net.Conn) error {
	toSend, err := fmtData(slc)
	fmt.Println("sent:", string(toSend))
	conn.Write(toSend)
	return err
}

//formats input as resp array of bulk strings
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
