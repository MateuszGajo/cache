package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Commands string

const (
	CLRF string = "\r\n"
	PING Commands = "PING"
	ECHO Commands = "ECHO"

)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	ln, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to run on port 6379")
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
	
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}

		go handleConenction(conn)
	}
	
}

func BuildResponse(message string) string {
	return fmt.Sprintf("$%v\r\n%s\r\n", len(message), message)
}



func handleConenction(conn net.Conn) {
	defer conn.Close();

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		input :=string(buf[:n])
		args := strings.Split(input, CLRF)
		if len(args) < 3 {
			fmt.Println("invalid command received:", input)
			return
		}

		command := args[2]

		var response string


		switch(command) {
		case "ping":
			response = "+PONG\r\n"
		case "echo":
			response = BuildResponse(args[4])
		default: {
			response = "-ERR unknown command\r\n"
			fmt.Println("invalid command received:", command)
		}
		}

		_, err = conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
		
	}
}
