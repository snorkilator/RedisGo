package main

import "testing"

var a [][]byte = [][]byte{[]byte("get"), []byte("a")}
var b []byte = []byte(`*2\r\n$3\r\nget\r\n$1\r\na\r\n`)

func TestFmtData(T *testing.T) {

}
