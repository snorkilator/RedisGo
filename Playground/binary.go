package main

import "fmt"

var a byte

var a map[[1000]byte][1000]byte

func main() {
	a += 255
	fmt.Println(string(a), a)
}
