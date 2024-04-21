package main

import (
	"fmt"
	"os"
	"strconv"
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

func pingCom (conn net.Conn){
	conn.Write([]byte("+PONG\r\n"))
}

func echoCom (conn net.Conn, input string){
	conn.Write([]byte("$"+ strconv.Itoa(len(input))+"\r\n"+input+"\r\n"))
}

// func BuildResponse(message string) string {
// 	return fmt.Sprintf("$%v\r\n%s\r\n", len(message), message)

// }


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


		switch(command) {
		case "ping":
			pingCom(conn)
		case "echo":
			echoCom(conn, args[4])
		default: {
			conn.Write([]byte("not found\r\n"))
		}
		}
		
	}
}
