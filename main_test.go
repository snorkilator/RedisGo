package main

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestFmtData(t *testing.T) {
	input := [][]byte{[]byte("get"), []byte("a")}

	want := []byte(`*2\r\n$3\r\nget\r\n$1\r\na\r\n`)
	got := fmtData(input)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", string(got), string(want))
	}
}

func TestHandleCommandError(t *testing.T) {
	// test error for incorrect command
	input := [][]byte{[]byte("net")}
	_, err := handleCommand(&input)
	if err.Error() != fmt.Sprintf("handleCommand: no such command %s, try set or get", input[0]) {
		t.Fatalf("handleCommand should throw error")
	}
}

func TestGetError(t *testing.T) {
	// error when there are too many
	tooManySlc := [][]byte{[]byte("get"), []byte("get"), []byte("get")}
	_, err := get(&tooManySlc)
	if err.Error() != "get: wrong number of arguments. want 2 but got 3" {
		t.Fatalf("get should return error for too many inputs")
	}

	tooFewSlc := [][]byte{[]byte("get")}
	_, err = get(&tooFewSlc)
	if err.Error() != "get: wrong number of arguments. want 2 but got 1" {
		t.Fatalf("get should return error for too few inputs")
	}

	// error when key is not found in database
	keyNotPresent := [][]byte{[]byte("get"), []byte("a")}
	_, err = get(&keyNotPresent)
	if err.Error() != fmt.Sprintf(`get: "%v" does not exist in database`, keyNotPresent[1]) {
		t.Fatalf("get should return error: no key found in database")
	}
}

func TestGetSetHappy(t *testing.T) {
	db = make(map[string][]byte)
	setA := [][]byte{[]byte("set"), []byte("a"), []byte("b")}
	getA := [][]byte{[]byte("get"), []byte("a")}
	expectOk := [][]byte{[]byte("OK")}

	ok, err := set(&setA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(ok, expectOk) {
		t.Fatalf(`expected "OK", got %v`, ok)
	}

	expected := [][]byte{[]byte("b")}
	result, err := get(&getA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("was not able to set and get")
	}

}

func TestSetError(t *testing.T) {
	// error when there are too many elements
	tooManySlc := [][]byte{[]byte("set"), []byte("set"), []byte("set"), []byte("set")}
	_, err := set(&tooManySlc)
	if err.Error() != "set: wrong number of arguments. want 3 but got 4" {
		t.Fatalf("get should return error if too many elements")
	}

	// error when there are too few elements
	tooFewSlc := [][]byte{[]byte("get")}
	_, err = set(&tooFewSlc)
	if err.Error() != "set: wrong number of arguments. want 3 but got 1" {
		t.Fatalf("get should return error for too few inputs")
	}

}

func TestParseHappy(t *testing.T) {
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
			t.Fatalf("unexpected error")
		}
		if !reflect.DeepEqual(got, e.expect) {
			t.Fatalf("got %v want %v", got, e.expect)
		}

	}
}

// start server and init database
// request connection and test for success
// listen for message from server in seperate thread
// send set command
// if listener doesn't receive response, error
// send get command, if not expected message, error
func TestServerAccept(t *testing.T) {
	go main()
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		t.Fatalf("error dialing")
	}
	defer conn.Close()

	_, err = conn.Write([]byte(`*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n`))
	if err != nil {
		t.Fatalf("error sending data")
	}
}
