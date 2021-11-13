package

// parser(d []byte) [][]byte
// parser stores each string from bulk string array in [][]byte
// 	checks that first byte is *
//  gets length of array
//		starting at 2nd byte, looks for byte sequence that represents \r\n
//		stores as int, 2nd byte to \r\n exclusive
//		