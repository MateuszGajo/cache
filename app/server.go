package main

import (
	"flag"
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
	INFO Commands = "INFO"
)

var CLRF string

func init() {
    os := runtime.GOOS

    if os == "windows" {
        CLRF = "\\r\\n" // fix this
    } else {
        CLRF = "\r\n"
    }
}

type CustomSetStore struct {
	Value      string
	ExpireAt  time.Time
}


var m = map[string]CustomSetStore{}
var lock = sync.RWMutex{}

type Replica struct {
	Port string
	Address string
}

var port int
var replica Replica

func init(){
	flag.IntVar(&port, "port", 6379, "port to listen to")
	flag.StringVar(&replica.Address, "replicaof", "", "master address")
	flag.Parse()
	if len(flag.Args()) > 0 {
		replica.Port = flag.Args()[0]
	}
}

type Server struct {
	role string
	replicaId string
	replicaOffSet int
}

// func NewServer(network string, port string) *Server {
// 	return &Server{
// 		address: fmt.Sprintf("0.0.0.0:%s", port),
// 		rep:     newReplication(master),
// 	}
// }

func handShake() error {
	conn, err := net.Dial("tcp", replica.Address + ":" + replica.Port)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", replica.Address, replica.Port)
	}

	defer conn.Close()

	conn.Write([]byte("*1\r\n$4\r\nping\r\n"))


	args := readInput(conn)

	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n" + strconv.Itoa((port)) + "\r\n"))
	args = readInput(conn)

	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
		

	return nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	serverCon := Server{
		role: "master",
	}
	if(replica != Replica {}) {
		serverCon.role = "slave"
		handShake()
	} else {
		serverCon.replicaOffSet = 0
		serverCon.replicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	}
	

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("Failed to run on port %d", port)
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
	
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}

		go handleConenction(conn, serverCon)
	}
	
}

func BuildResponse(message string) string {
 return fmt.Sprintf("$%v\r\n%s\r\n", len(message), message)
}

func BuildResponses(messages []string) string {
	res := fmt.Sprintf("*%v\r\n", len(messages))

	for _, val := range messages {
		res += BuildResponse(val)
	}

	return res 
}

func handleRemove(key string) {
	lock.Lock()
	defer lock.Unlock()
	delete(m, key) 
	//TO DO
} 

func handleSet(key, value string, expiryTime *int) string {
	lock.Lock()
	defer lock.Unlock()
	if expiryTime != nil {
		m[key] = CustomSetStore {
			Value:  value,
			ExpireAt: time.Now().Add(time.Duration(*expiryTime) * time.Millisecond),
		}
	} else {
		m[key] = CustomSetStore {
			Value:  value,
			ExpireAt: time.Time{},
		}
	}
	
	
	return "+OK\r\n"
}

func handleGet(key string) string {
	defer lock.RUnlock()
	lock.RLock()
	r, ok := m[key]
	if !ok || (time.Now().After(r.ExpireAt) && r.ExpireAt != time.Time{}) {
		return "$-1\r\n"
	}
	
	return BuildResponse(r.Value)
}

func RESPSimpleString (input string) string{
	arg := strings.Split(input, CLRF)[0]
	arg = strings.Trim(arg, "+")

	return arg
}

func RESPArray(input string) [] string {
	fmt.Printf("\n %v \n", input)
	args := strings.Split(input, CLRF)

	fmt.Println(args)

	command  := make([]string, 0, ((len(args) -1) /2))

	for i :=2; i< len(args); i = i+2 {
		command = append(command, args[i])
	}

	return command
}

func RESPBulkString(input string) []string {
	fmt.Printf("\n %v \n", input)
	args := strings.Split(input, CLRF)

	fmt.Println(args)

	command  := make([]string, 0, ((len(args) -1) /2))

	for i :=2; i< len(args); i = i+2 {
		command = append(command, args[i])
	}

	return command
}

func readInput(conn net.Conn) []string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Print("problem reading", err)
		os.Exit(1)
	}

	command := []string {}
	input :=string(buf[:n])
	fmt.Println(input)
	switch(input[0]) {
	case 43:
		command = append(command,RESPSimpleString(input))
	case 36:
		command = RESPBulkString(input)
	case 42:
		command = RESPArray(input)
	default:
		fmt.Print("Unknown RESP enconfing")
		os.Exit(1)
	}
	fmt.Println(command)

	return command

}

func handleConenction(conn net.Conn, serverCon Server) {
	defer conn.Close();

	for {
		args := readInput(conn)
		fmt.Print(args)

		if len(args) == 0 {
			fmt.Println("No argument passed")
			os.Exit(1)
		}



		command := Commands(strings.ToUpper(args[0]))

		var response string

		switch(command) {
		case PING:
			response = "+PONG\r\n"
		case ECHO:
			response = BuildResponse(args[4])
		case SET:	
			switch(len(args)){
			case 4:	
				// command := args[3]
				timeMs, err := strconv.Atoi(args[4])
				if err != nil {
					fmt.Print("invalid time")
					os.Exit(1)
				}
				response = handleSet(args[1], args[2], &timeMs)
			case 3:
				response = handleSet(args[1], args[2], nil)
			}	
		case GET:
			response = handleGet(args[1])
		case INFO:
			response = BuildResponse("role:"+ serverCon.role +"master_replid:" + serverCon.replicaId +"master_repl_offset:" + strconv.Itoa(serverCon.replicaOffSet))
			fmt.Printf("response, %v",response)
		default: {
			response = "-ERR unknown command\r\n"
			fmt.Println("invalid command received:", command)
		}
		}
		

		_, err := conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
		
	}
}
