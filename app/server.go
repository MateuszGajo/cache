package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Commands string

const (
	PING Commands = "PING"
	ECHO Commands = "ECHO"
	SET Commands = "SET"
	GET Commands = "GET"
	INFO Commands = "INFO"
	REPLCONF Commands = "REPLCONF"
	PSYNC Commands = "PSYNC"
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


func handShake() net.Conn {
	conn, err := net.Dial("tcp", replica.Address + ":" + replica.Port)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", replica.Address, replica.Port)
	}

	defer conn.Close()

	conn.Write([]byte("*1\r\n$4\r\nping\r\n"))


	inputComm := readInput(conn)
	args := inputComm.commandStr

	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n" + strconv.Itoa((port)) + "\r\n")) // lets build it
	inputComm = readInput(conn)
	args = inputComm.commandStr

	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n")) // same lets build it
	inputComm = readInput(conn)
	args = inputComm.commandStr
	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
		

	return conn
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	serverCon := Server{
		role: "master",
	}
	
	

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("Failed to run on port %d", port)
		os.Exit(1)
	}

	defer func(){
		ln.Close() // is it closed?
		fmt.Print("close")
	}()

	for {
		conn, err := ln.Accept()
		var replConn net.Conn
		if(replica != Replica {}) {
			serverCon.role = "slave"
			replConn = handShake()
		} else {
			serverCon.replicaOffSet = 0
			serverCon.replicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		}
	
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}

		go handleConenction(conn, serverCon)
	}
	
}

type CommandInput struct{
	commandStr []string
	commandByte string
}

func readInput(conn net.Conn) CommandInput{
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Print("problem reading", err)
		os.Exit(1)
	}

	command := []string {}
	input :=string(buf[:n])

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

	return CommandInput{
		commandStr: command,
		commandByte: input,
	}

}

func propagte (command []byte, conn net.Conn) {
	fmt.Print("enter??")
	fmt.Print(command)
	_, err := conn.Write(command)

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}
}

func handleConenction(conn net.Conn, serverCon Server) {
	defer func() {
		conn.Close(); // is it closed?
		fmt.Print("close")
	}()
	//we need to establish connection with replica

	var replConn net.Conn

	for {
		comamndInput := readInput(conn)
		args := comamndInput.commandStr
		fmt.Print(args)

		if len(args) == 0 {
			fmt.Println("No argument passed")
			os.Exit(1)
		}

		command := Commands(strings.ToUpper(args[0]))

		var response string

		switch(command) {
		case PING: // fix failing ping problem
			response = Ping()
		case ECHO:
			response = Echo(args)
		case SET:	
			response = Set(args)
			fmt.Print(comamndInput)
			if replConn != nil {
			propagte([]byte(comamndInput.commandByte), replConn)
			}
		case GET:
			response = Get(args)
			if replConn != nil {
				propagte([]byte(comamndInput.commandByte), replConn)
			}
		case INFO:
			response = Info(args, serverCon)
			replConn = conn
		case REPLCONF:
			response = ReplConf()
		case PSYNC: 
			response = Psync(conn, serverCon)
			continue // we need to resolve this problem,
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
