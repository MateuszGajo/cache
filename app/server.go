package main

import (
	"errors"
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

// func init(){
	
// }

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


func handShake(){
	fmt.Print("replica handshake?")
	conn, err := net.Dial("tcp", replica.Address + ":" + replica.Port)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", replica.Address, replica.Port)
	}

	conn.Write([]byte("*1"+CLRF+"$4"+CLRF+"ping"+CLRF))


	inputComm := readInput(conn)
	args := inputComm.commandStr

	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$8"+CLRF+"REPLCONF"+CLRF+"$14"+CLRF+"listening-port"+CLRF+"$4"+CLRF + strconv.Itoa((port)) + CLRF)) // lets build it
	inputComm = readInput(conn)
	args = inputComm.commandStr

	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$8"+CLRF+"REPLCONF"+CLRF+"$4"+CLRF+"capa"+CLRF+"$6"+CLRF+"psync2"+CLRF)) // same lets build it
	inputComm = readInput(conn)
	args = inputComm.commandStr
	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$5"+CLRF+"PSYNC"+CLRF+"$1"+CLRF+"?"+CLRF+"$2"+CLRF+"-1"+CLRF))
	fmt.Println("sending psync")
	inputComm = readInput(conn)
	args = inputComm.commandStr

	fmt.Print("response")
	fmt.Println(inputComm.commandStr)
	inputComm = readInput(conn)
	fmt.Print("i hope we are not here")



	// if args[0] != "FULLRESYNC" {
	// 	fmt.Println("hello")
	// 	fmt.Println(args[0])
	// 	fmt.Print("Response its invalid")
	// 	os.Exit(1)
	// }

	// defer conn.Close()

	go handleConenction(conn, Server{})
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	var aa int

	flag.IntVar(&port, "port", 6379, "port to listen to")
	flag.StringVar(&replica.Address, "replicaof", "", "master address")
	flag.Parse()
	fmt.Println("port", port)
		aa = port

	if len(flag.Args()) > 0 {
		
		fmt.Println(flag.Args())
		replica.Port = flag.Args()[0]

	}

	// Uncomment this block to pass the first stage
	serverCon := Server{
		role: "master",
	}
		fmt.Print(fmt.Print(aa))
	

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", aa))
	fmt.Print(fmt.Printf("127.0.0.1:%d", aa))

	if err != nil {
		fmt.Printf("Failed to run on port %d", port)
		os.Exit(1)
	}

	defer func(){
		ln.Close() // is it closed?
		fmt.Print("close")
	}()
	fmt.Print(replica)
	fmt.Print(replica != Replica{})

	for {
		fmt.Print("ddds1?")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}
		fmt.Print("ddds?")
		if(replica != Replica{}) {
			fmt.Println("go for handshake??")
			serverCon.role = "slave"
			handShake()
		} else {
			fmt.Println("halo")
			serverCon.replicaOffSet = 0
			serverCon.replicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		}
	
	

		go handleConenction(conn, serverCon)
	}
	
}

type CommandInput struct{
	commandStr []string
	commandByte string
}
// type MyConn struct {
// 	net.Conn
// 	ignoreWrites bool
// }

func readInput(conn net.Conn) CommandInput{
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Print("problem reading", err)
		os.Exit(1)
	}

	command := []string {}
	input :=string(buf[:n])


fmt.Println("hello")
fmt.Println(input)
fmt.Println(input[0])
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


	return CommandInput{
		commandStr: command,
		commandByte: input,
	}

}
var replConn net.Conn

func propagte (conn net.Conn, command []byte) {
	fmt.Println("propagte")
	if replConn == nil {
		return
	}
	replConn.Write(command)
}

// func (conn MyConn) Write(b []byte) (n int, err error) {
// 	if conn.ignoreWrites {
// 		return len(b), nil
// 	}
// 	fmt.Println("write")
// 	fmt.Println(b)
// 	return conn.Conn.Write(b)
// }

var whitelistProp = map[Commands]bool{"SET": true}

func handleConenction(conn net.Conn, serverCon Server) {
	defer func() {
		conn.Close(); // is it closed?
		fmt.Print("close")
	}()

	for {
		fmt.Println("read another")
		comamndInput := readInput(conn)
		args := comamndInput.commandStr

		if len(args) == 0 {
			fmt.Println("No argument passed")
			os.Exit(1)
		}

		command := Commands(strings.ToUpper(args[0]))

		var err error

		// connectionFromMaster := strings.Contains(conn.RemoteAddr().String(), "6379") // fix port
		// if(connectionFromMaster) {
		// 	conn = MyConn{
		// 		ignoreWrites: true,
		// 		Conn: conn,
		// 	}
		// } else {
		// 	conn = MyConn{
		// 		ignoreWrites: false,
		// 		Conn: conn,
		// 	}
		// } // is it to much coping??

		if whitelistProp[command] && replConn != nil {
			propagte(conn, []byte(comamndInput.commandByte))
		}

		switch(command) {
		case PING:
			fmt.Println("reading ping")
			err = Ping(conn,) // test if reciver does not breaking something
		case ECHO:
			fmt.Println("reading echi")
			err = Echo(conn,args)
		case SET:	
		fmt.Println("reading set")
			err = Set(conn,args)
		case GET:
			fmt.Println("reading get")
			err = Get(conn,args)
		case INFO:
			fmt.Println("reading info")
			err = Info(conn,args, serverCon)
		case REPLCONF:
			fmt.Println("reading replconfg")
			err = ReplConf(conn)
			replConn = conn
		case PSYNC: 
			fmt.Println("reading psync")
			err = Psync(conn,serverCon)
			replConn = conn
		default: {
			err = errors.New(fmt.Sprintf("invalid command received:%v", command))
			fmt.Println("invalid command received:", command)
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
		}
		
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
	}
}
