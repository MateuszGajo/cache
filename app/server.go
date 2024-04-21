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

func parseResp (data string) []string {
	result := strings.Split(data, "$")

	if len(result) > 1 {
		result = result[1:]
		fmt.Println(result)
		for i := range result {
			parts := strings.Split(result[i], "\\r\\n")
			if(len(parts) > 1) {
				result[i] = parts[1] 
			} else {
				result[i] = parts[0] 
			}
			
		}
	}

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
		var args []string;
		fmt.Print(input)
		fmt.Println(input[0])
		if(input[0] == 42) {
			args = parseResp(input)
		} else {
			args = strings.Split(input, " ")
		}
		fmt.Println(args)
		command := args[1]
		values := args[2:]


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
