package main

import (
	"fmt"
	"os"

	// Uncomment this block to pass the first stage
	"net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to run on port 6379")
		os.Exit(1)
	}
	_, err = ln.Accept()
	ln.Accept()

	if err != nil {
		fmt.Println("Error accepting connection", err.Error())
		os.Exit(1)
	}
}
