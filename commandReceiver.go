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

func responder(s string) error {
	fmt.Println(s)
	return nil
}

var input = `*4\r\n$3\r\nset\r\n$1\r\na\r\n$2\r\n21\r\n$2\r\ner\r\n`

//Takes array of bulk string, and outputs a slice of strings containing the elements of the array
func parser(s string) ([]string, error) {
	var parsed []string

	if s[0] != byte('*') {
		return nil, fmt.Errorf("input is not Array of Bulk Strings type: want *, got %v", string(s[0]))
	}

	arrayLen, currentIndex, err := getLen(1, &s)
	if err != nil {
		return nil, err
	}

	for i := 0; i < arrayLen; i, currentIndex = i+1, currentIndex+4 {

		if s[currentIndex] != '$' {
			return nil, fmt.Errorf("Expected type symbol '$' for element %v in array but got %v at index %v", i+1, string(s[currentIndex]), currentIndex)
		}
		currentIndex++

		var strLen int
		strLen, currentIndex, err = getLen(currentIndex, &s)
		if err != nil {
			return nil, err
		}

		tempStr := s[currentIndex : currentIndex+strLen]
		currentIndex += strLen
		parsed = append(parsed, tempStr)
		if s[currentIndex:currentIndex+4] != `\r\n` {
			return nil, fmt.Errorf(`expected \r\n after %v but found %v starting at index %v`, parsed, s[currentIndex:currentIndex+4], currentIndex)
		}
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
		if (*s)[i:i+4] == `\r\n` {
			found = true
			endofNum = i - 1
		}
	}
	lenS := (*s)[beginStr : endofNum+1]

	len, err := strconv.Atoi(lenS)
	if err != nil {
		return 0, 0, fmt.Errorf("Array length is not valid number: want number but got %v", lenS)
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
	// s, _ := receiver()
	// fmt.Printf("%v", s)
	input, err := receiver()
	if err != nil {
		fmt.Println("receiver error:", err)
	}

	s, err := parser(input)
	if err != nil {
		fmt.Println(err)
	}

	err = commandHandler(&s)
	if err != nil {
		fmt.Println("commanderHandler():", err)
	}
}
