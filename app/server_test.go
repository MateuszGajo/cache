package main

import (
	"flag"
	"fmt"
	"net"
	"testing"
)

var _ = func() bool {
    testing.Init()
	// main()
    return true
}()

func TestGGGGG(t *testing.T) {
	go main()
	flag.Set("port", "8888")
	// flag.Parse()
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
    if err != nil {
        t.Error("could not connect to server: ", err)
    }
    defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF, CLRF, CLRF, CLRF, CLRF, CLRF, CLRF)))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		 t.Error("could not connect to server: ", err)
	}

	input :=string(buf[:n])

	fmt.Print(input)
	if input != "+OK" + CLRF {
		t.Errorf("Wront output, expected %v, got %v ","+Ok" + CLRF, input)
	}

}