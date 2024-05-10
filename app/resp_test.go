package main

import (
	"flag"
	"testing"
)

// func init() {
//     flag.StringVar(&foo, "foo", "", "the foo bar bang")
//     flag.Parse()
// }

var _ = func() bool {
    testing.Init()
    return true
}()

func TestSimpleString(t *testing.T) {
	flag.Set("port", "8888")
	// flag.Parse()
	resp := BuildSimpleString("test")

	if resp != "+test" + CLRF {
		t.Fatalf("Expected:%v, got:%v", "+test"+ CLRF, resp)
	}
}