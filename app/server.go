package main

import (
	"fmt"
	"os"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	// "os"
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

func pingCom (conn net.Conn, args ...string){
	conn.Write([]byte("+PONG\r\n"))
}

func echoCom (conn net.Conn, args ...string){
	fmt.Println(args)
	input := args[0]
	conn.Write([]byte("$"+string(len(input))+"\r\n"+input+"r\n"))
}


type Commands string

const (
	CLRF string = "\r\n"
	PING Commands = "PING"
	ECHO Commands = "ECHO"
)

func parseResp (data string) []string {
	result := strings.Split(data, CLRF)

	return result
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
		args := parseResp(input)
		fmt.Print(args)
		if len(args) < 3 {
			fmt.Println("invalid command received:", input)
			return
		}

		fmt.Println(args)
		command := args[2]
		values := args[3:]


		switch(command) {
		case "ping":
			pingCom(conn, values...)
		case "echo":
			echoCom(conn, values...)
		default: {
			conn.Write([]byte("not found\r\n"))
		}
		}
		
	}
}
