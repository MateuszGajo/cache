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
	resp := BuildRESPArray([]string{"REPLCONF","ACK", "0"})

	if resp != fmt.Sprintf("*3%v$8%vREPLCONF%v$3%vACK%v$1%v0%v", CLRF,CLRF,CLRF,CLRF,CLRF ,CLRF,CLRF) {
		t.Fatalf("Expected:%v, got:%v", fmt.Sprintf("*3%v$8%vREPLCONF%v$3%vACK%v$1%v0%v", CLRF,CLRF,CLRF,CLRF,CLRF ,CLRF,CLRF), resp)
	}
}