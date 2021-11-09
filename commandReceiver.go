package main

import (
	"bufio"
	"fmt"
	"os"
)

func receiver() (s string, err error) {
	reader := bufio.NewReader(os.Stdin)
	s, err = reader.ReadString('\n')
	return
}

var input = `*3\r\n$3\r\nset\r\n$1\r\na\r\n$2\r\n21\r\n`

func parser(s string) (a []string, err error) {
	if s[0] != byte('*') {
		return nil, fmt.Errorf("input is not Array of Bulk Strings type: want *, got %v", string(s[0]))
	}

	var endofNum int
	for i, found := 1, false; !found; i++ {
		if s[i:i+4] == `\r\n` {
			found = true
			endofNum = i - 1
		}
	}
	fmt.Println(s[1 : endofNum+1])
	//find length of bulk string array
	//number is stored between * idex 0 and first \r\n
	return nil, nil
}

//parser
//		decodes RESP array, and stores in slice
//			check datatype is orrect (*)
//			read length of array
//				read characters from index 1 onward
//				until you reach \r\n
// 				numerical length of array is index 1 to \r\n exclusive
//			store length of array
//			initialize array of strings of length of array
//			find each element in input array and store in new array
//				start after first \r\n
//				check if datatype is correct ($)
//				find length of string
//					find everythig between $ and \r\n
//				store as length of string
//				extract string
//					find first \r\n
//					grab length of string starting after \r\n
//				take string and store in first empty space in array
//
func main() {
	// s, _ := receiver()
	// fmt.Printf("%v", s)
	parser(input)
}

/*
command: set a 21
(set key `a` to value `21`)
RESP:
*3\r\n$3\r\nset\r\n$1\r\na\r\n$2\r\n21\r\n
*/
