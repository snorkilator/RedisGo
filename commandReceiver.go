package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//receives command from client and returns it as a string
func receiver() (s string, err error) {
	reader := bufio.NewReader(os.Stdin)
	s, err = reader.ReadString('\n')
	return
}

//sends information back to client
func responder(s string) error {
	fmt.Println(s)
	return nil
}

var input []byte = []byte(`*3\r\n$3\r\nset\r\n$1\r\na\r\n$2\r\n23\r\n`)

//Takes array of bulk string, and outputs a slice of strings containing the elements of the array
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
			return nil, fmt.Errorf("Parser: Expected type symbol '$' for element %v in array but got %v at index %v", i+1, string(s[currentIndex]), currentIndex)
		}
		currentIndex++

		var strLen int
		strLen, currentIndex, err = getLen(currentIndex, &s)
		if err != nil {
			return nil, err
		}

		if currentIndex+strLen > len(s)-1 {
			return nil, fmt.Errorf(`Parser: %v indexed element does not exist in array`, i+1)
		}
		tempStr := b[currentIndex : currentIndex+strLen]
		currentIndex += strLen

		if len(s)-1 < currentIndex+3 {
			return nil, fmt.Errorf(`Parser: array element %v does not terminate with \r\n`, i)
		}
		if s[currentIndex:currentIndex+4] != `\r\n` {
			return nil, fmt.Errorf(`Parser: expected \r\n but found %v starting at index %v`, s[currentIndex:currentIndex+4], currentIndex)
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

func commandHandler(slc *[]string) error {
	switch (*slc)[0] {
	case "set":
		set(slc)
	case "get":
		get(slc)
	default:
		responder("No such command " + (*slc)[0] + ", try SET or GET")
	}

	return nil
}

func set(*[]string) {
	// check that there are three elements
	//		throw error if false
	//
}
func get(*[]string) {}

func main() {

	// input, err := receiver()
	// if err != nil {
	// 	fmt.Println("receiver error:", err)
	// }

	s, err := parser(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	// err = commandHandler(&s)
	// if err != nil {
	// 	fmt.Println("commanderHandler():", err)
	// }
}
