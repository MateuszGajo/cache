package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Commands string

const (
	PING Commands = "PING"
	ECHO Commands = "ECHO"
	SET Commands = "SET"
	GET Commands = "GET"
)

var CLRF string

func init() {
    os := runtime.GOOS

    if os == "windows" {
        CLRF = "\\r\\n"
    } else {
        CLRF = "\r\n"
    }
}


var m = map[string]string{}
var lock = sync.RWMutex{}

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

func handleRemove(key string) {
	lock.Lock()
	defer lock.Unlock()
	delete(m, key)
} 

func handleSet(key, value string, triggerMs *int) string {
	lock.Lock()
	defer lock.Unlock()
	m[key] = value
	if triggerMs != nil {
		time.AfterFunc(time.Duration(*triggerMs) *time.Millisecond, func (){handleRemove(key)})
	}
	return "+OK\r\n"
}

func handleGet(key string) string {
	defer lock.RUnlock()
	lock.RLock()
	val, ok := m[key]
	if !ok {
		return "$-1\r\n"
	}
	
	return BuildResponse(val)
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
		fmt.Print(args)
		if len(args) < 3 {
			fmt.Println("invalid command received:", input)
			return
		}

		command := Commands(strings.ToUpper(args[2]))

		var response string

		
		

		switch(command) {
		case PING:
			response = "+PONG\r\n"
		case ECHO:
			response = BuildResponse(args[4])
		case SET:	
			switch(len(args)){
			case 9:	
				timeMs, err := strconv.Atoi(args[8])
				if err != nil {
					fmt.Print("invalid time")
					os.Exit(1)
				}
				response = handleSet(args[4], args[6], &timeMs)
			case 7:
				response = handleSet(args[4], args[6], nil)
			}	
		case GET:
			response = handleGet(args[4])
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
