package main

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"redis/server"
	"reflect"
	"testing"
	"time"
)

func TestFmtData(t *testing.T) {
	input := [][]byte{[]byte("get"), []byte("a")}

	want := []byte("*2\r\n$3\r\nget\r\n$1\r\na\r\n")
	got := fmtData(input)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v,\nwant %v \n", got, want)
	}
}

func TestHandleCommandError(t *testing.T) {
	// test error for incorrect command
	input := [][]byte{[]byte("net")}
	_, err := handleCommand(&input)
	if err.Error() != fmt.Sprintf("handleCommand: no such command %s, try set or get", input[0]) {
		t.Fatalf("handleCommand should throw error %v", err)
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

func TestGetLen(t *testing.T) {
	input := []byte("12\r\n132\r\n")
	want := 132
	lenB, begin, err := getLen(4, &input)
	if lenB != want {
		t.Fatalf("got %d, want %d", lenB, want)
	}
	if begin != len(input) {
		t.Fatalf("got %d, want %d", begin, len(input))
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestParseHappy(t *testing.T) {
	tc := []struct {
		input  []byte
		expect [][]byte
	}{
		{
			input:  []byte("*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n"),
			expect: [][]byte{[]byte("set"), []byte("a"), []byte("b")},
		},
		{
			input:  []byte("*2\r\n$3\r\nget\r\n$1\r\na\r\n"),
			expect: [][]byte{[]byte("get"), []byte("a")},
		},
	}

	for _, c := range tc {
		got, err := parse(c.input)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if !reflect.DeepEqual(got, c.expect) {
			t.Fatalf("got %v want %v", got, c.expect)
		}

	}
}

func TestServer(t *testing.T) {

	// init server and connect to it
	go main()
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer conn.Close()
	defer server.Server.Close()

	t.Run("TestSendErr", func(t *testing.T) {
		_, err := conn.Write([]byte("*2\r\n$3\r\ngt\r\n$1\r\na\r\n"))
		if err != nil {
			t.Fatalf("error sending command: %v", err)
		}

		data := make([]byte, 512)
		_, err = conn.Read(data)
		if err != nil {
			t.Fatalf("error reading response: %v", err)
		}
		fmt.Printf("%s", data)
		if data[0] != '-' {
			t.Fatalf("did not receive error message")
		}

	})
	t.Run("TestSetCommand", func(t *testing.T) {
		_, err = conn.Write([]byte("*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n"))
		if err != nil {
			t.Fatalf("error sending data: %v", err)
		}

		want := []byte("*1\r\n$2\r\nOK\r\n")
		got := make([]byte, len(want)) //don't use length of want

		_, err = conn.Read(got) //add a trimmer to got to trim zeros
		if err != nil {
			t.Fatalf("error reading set response")
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("TestGetCommand (Get must work first)", func(t *testing.T) {
		_, err = conn.Write([]byte("*2\r\n$3\r\nget\r\n$1\r\na\r\n"))
		if err != nil {
			t.Fatalf("error sending data: %v", err)
		}

		want := []byte("*1\r\n$1\r\nb\r\n")
		got := make([]byte, len(want))

		_, err = conn.Read(got)
		if err != nil {
			t.Fatalf("error reading get response: %v", err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
	t.Run("redis-cli", func(t *testing.T) {
		cmd := exec.Command("redis-cli", "-h", server.CONN_HOST, "-p", server.CONN_PORT)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}

		err = cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		_, err = stdin.Write([]byte("set g qwerty\n"))
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 7)
		_, err = stdout.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf[:3], []byte{79, 75, 10}) {
			t.Fatalf("was not able to perform set command using redis-cli")
		}
		_, err = stdin.Write([]byte("get g\n"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = stdout.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf, []byte{113, 119, 101, 114, 116, 121, 10}) {
			t.Fatalf("was not able to perform get command using redis-cli")
		}
	})
}
