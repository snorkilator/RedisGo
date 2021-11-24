package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFmtData(T *testing.T) {
	var input [][]byte = [][]byte{[]byte("get"), []byte("a")}

	var want []byte = []byte(`*2\r\n$3\r\nget\r\n$1\r\na\r\n`)
	got, _ := fmtData(input)

	if !reflect.DeepEqual(got, want) {
		T.Fatalf("got %v, want %v", string(got), string(want))
	}
}

func TestHandleCommandError(T *testing.T) {
	// test error for incorrect command
	input := [][]byte{[]byte("net")}
	_, err := handleCommand(&input)
	if err.Error() != fmt.Sprintf("handleCommand: no such command %s, try set or get", input[0]) {
		T.Fatalf("handleCommand should throw error")
	}
}

func TestGetError(T *testing.T) {
	// error when there are too many
	tooManySlc := [][]byte{[]byte("get"), []byte("get"), []byte("get")}
	_, err := get(&tooManySlc)
	if err.Error() != "get: wrong number of arguments. want 2 but got 3" {
		T.Fatalf("get should return error for too many inputs")
	}

	tooFewSlc := [][]byte{[]byte("get")}
	_, err = get(&tooFewSlc)
	if err.Error() != "get: wrong number of arguments. want 2 but got 1" {
		T.Fatalf("get should return error for too few inputs")
	}

	// error when key is not found in database
	keyNotPresent := [][]byte{[]byte("get"), []byte("a")}
	_, err = get(&keyNotPresent)
	if err.Error() != fmt.Sprintf(`get: "%v" does not exist in database`, keyNotPresent[1]) {
		T.Fatalf("get should return error: no key found in database")
	}
}

func TestGetSetHappy(T *testing.T) {
	db = make(map[string][]byte)
	setA := [][]byte{[]byte("set"), []byte("a"), []byte("b")}
	getA := [][]byte{[]byte("get"), []byte("a")}
	expectOk := [][]byte{[]byte("OK")}

	ok, err := set(&setA)
	if err != nil {
		T.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(ok, expectOk) {
		T.Fatalf(`expected "OK", got %v`, ok)
	}

	expected := [][]byte{[]byte("b")}
	result, err := get(&getA)
	if err != nil {
		T.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		T.Fatalf("was not able to set and get")
	}

}

func TestSetError(T *testing.T) {
	// error when there are too many elements
	tooManySlc := [][]byte{[]byte("set"), []byte("set"), []byte("set"), []byte("set")}
	_, err := set(&tooManySlc)
	if err.Error() != "set: wrong number of arguments. want 3 but got 4" {
		T.Fatalf("get should return error if too many elements")
	}

	// error when there are too few elements
	tooFewSlc := [][]byte{[]byte("get")}
	_, err = set(&tooFewSlc)
	if err.Error() != "set: wrong number of arguments. want 3 but got 1" {
		T.Fatalf("get should return error for too few inputs")
	}

}

func TestParseHappy(T *testing.T) {
	tc := []struct {
		input  []byte
		expect [][]byte
	}{
		{
			input:  []byte(`*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n`),
			expect: [][]byte{[]byte("set"), []byte("a"), []byte("b")},
		},
		{
			input:  []byte(`*2\r\n$3\r\nget\r\n$1\r\na\r\n`),
			expect: [][]byte{[]byte("get"), []byte("a")},
		},
	}

	for _, e := range tc {
		got, err := parse(e.input)
		if err != nil {
			T.Fatalf("unexpected error")
		}
		if !reflect.DeepEqual(got, e.expect) {
			T.Fatalf("got %v want %v", got, e.expect)
		}

	}
}
