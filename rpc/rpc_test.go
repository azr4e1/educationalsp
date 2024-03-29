package rpc_test

import (
	"fmt"
	"testing"

	"github.com/azr4e1/educationalsp/rpc"
)

type EncodingExample struct {
	Method string
}

func TestEncode(t *testing.T) {
	input := EncodingExample{Method: "hello"}
	want := fmt.Sprintf("Content-Length: 18\r\n\r\n{\"Method\":\"hello\"}")
	got := rpc.EncodeMessage(input)

	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestDecode(t *testing.T) {
	input := []byte(fmt.Sprintf("Content-Length: 18\r\n\r\n{\"Method\":\"hello\"}"))
	methodWanted, lengthWanted := "hello", 18
	gotMethod, gotContent, err := rpc.DecodeMessage(input)
	gotLength := len(gotContent)
	if err != nil {
		t.Fatal(err)
	}

	if gotLength != lengthWanted {
		t.Errorf("want %d, got %d", lengthWanted, gotLength)
	}

	if gotMethod != methodWanted {
		t.Errorf("want %s, got %s", methodWanted, gotMethod)
	}
}
