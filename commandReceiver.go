package main

import (
	"bufio"
	"fmt"
	"os"
)

func receiver() (s string, err error) {
	reader := bufio.NewReader(os.Stdin)
	s, err = reader.ReadString('\n')
	s, err = reader.ReadString('\n')
	return
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
//			find each element in array and store in array
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
	s, _ := receiver()
	fmt.Printf("%v", s)
}

/*
command: set a 21
(set key `a` to value `21`)
RESP:
*3\r\n$3\r\nset\r\n$1\r\na\r\n$2\r\n21\r\n
*/
