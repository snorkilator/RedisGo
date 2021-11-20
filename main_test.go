package main

import (
	"fmt"
	"reflect"
	"testing"
)

var a [][]byte = [][]byte{[]byte("get"), []byte("a")}

var B []byte = []byte(`*2\r\n$3\r\nget\r\n$1\r\na\r\n`)

func TestFmtData(T *testing.T) {
	r, _ := fmtData(a)

	if !reflect.DeepEqual(r, B) {
		T.Fatalf("got %v, want %v", string(r), string(B))
	}
}

func TestHandleCommand(T *testing.T) {
	input := [][]byte{[]byte("net")}
	_, err := handleCommand(&input)
	if err == nil {
		T.Fatalf("handleCommand should throw error")
	}
	fmt.Println(err)
}
