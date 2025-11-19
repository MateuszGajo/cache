package main

import (
	"fmt"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestSimpleString(t *testing.T) {
	resp := BuildPrimitiveRESPArray([]string{"REPLCONF", "ACK", "0"})

	expectedOutput := fmt.Sprintf("*3%v$8%vREPLCONF%v$3%vACK%v$1%v0%v", CLRF, CLRF, CLRF, CLRF, CLRF, CLRF, CLRF)

	if resp != expectedOutput {
		t.Fatalf("Expected:%v, got:%v", expectedOutput, resp)
	}
}

func TestSimpleError(t *testing.T) {
	resp := BuildSimpleErrorWithErrType("Error", "This is error")

	if resp != "-Error This is error\r\n" {
		t.Fatalf("Expected:%v, got:%q", "-Error This is error", resp)
	}
}
