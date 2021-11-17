package main

import (
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
