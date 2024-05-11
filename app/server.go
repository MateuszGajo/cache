package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
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
		fmt.Println("port", port)
		fmt.Println(flag.Args())
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


func handShake(){
	fmt.Print("replica handshake?")
	fmt.Println("tcp", replica.Address + ":" + replica.Port)
	conn, err := net.Dial("tcp", replica.Address + ":" + replica.Port)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", replica.Address, replica.Port)
	}

	_, err = conn.Write([]byte("*1"+CLRF+"$4"+CLRF+"ping"+CLRF))

	if err != nil {
		fmt.Print("1 error while pinging master replica", err)
		conn.Close()
		return
	}


	inputComm, err := readInput(conn)
	if err != nil {
		fmt.Print("error while pinging master replica", err)
		conn.Close()
		return
	}
	args := inputComm.commandStr[0]

	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$8"+CLRF+"REPLCONF"+CLRF+"$14"+CLRF+"listening-port"+CLRF+"$4"+CLRF + strconv.Itoa((port)) + CLRF)) // lets build it
	inputComm, err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf listening-port master replica")
		conn.Close()
		return
	}
	args = inputComm.commandStr[0]

	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$8"+CLRF+"REPLCONF"+CLRF+"$4"+CLRF+"capa"+CLRF+"$6"+CLRF+"psync2"+CLRF)) // same lets build it
	inputComm,err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf capa  master replica")
		conn.Close()
		return
	}
	args = inputComm.commandStr[0]
	if args[0] != "OK" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte("*3"+CLRF+"$5"+CLRF+"PSYNC"+CLRF+"$1"+CLRF+"?"+CLRF+"$2"+CLRF+"-1"+CLRF))
	fmt.Println("sending psync")
	inputComm, err = readInput(conn)
	if err != nil {
		fmt.Print("error while psync master replica")
		conn.Close()
		return
	}
	// args = inputComm.commandStr[0]

	fmt.Print("response")
	fmt.Println(inputComm.commandStr)
	inputComm,err = readInput(conn)
	if err != nil {
		fmt.Print("error while psync second read replica")
		conn.Close()
		return
	}
	// fmt.Print("after handshake....")
	// fmt.Print(inputComm.commandStr)



	// if args[0] != "FULLRESYNC" {
	// 	fmt.Println("hello")
	// 	fmt.Println(args[0])
	// 	fmt.Print("Response its invalid")
	// 	os.Exit(1)
	// }

	// defer conn.Close()

	// return nil;

	go handleConenction(MyConn{Conn: conn, ignoreWrites: false, ID: strconv.Itoa(rand.IntN(100))}, Server{}) 
}

func main() {
	fmt.Print(os.Args)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	serverCon := Server{
		role: "master",
	}
		fmt.Print(fmt.Print(port))
	

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	fmt.Print(fmt.Printf("127.0.0.1:%d", port))

	if err != nil {
		fmt.Printf("Failed to run on port %d", port)
		os.Exit(1)
	}

	defer func(){
		fmt.Print("close")
		ln.Close() // is it closed?
	}()
	fmt.Print(replica)
	fmt.Print(replica != Replica{})

	if(replica != Replica{}) {
		fmt.Print("replica")
		fmt.Println("go for handshake??")
		serverCon.role = "slave"
		handShake()
	} else {
		fmt.Println("halo")
		serverCon.replicaOffSet = 0
		serverCon.replicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	}


	for {
		fmt.Print("ddds1?")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}
		fmt.Print("ddds?")
	

		go handleConenction(MyConn{Conn: conn, ignoreWrites: false, ID: strconv.Itoa(rand.IntN(100)) }, serverCon)
	}
	
}

type CommandInput struct{
	commandStr [][]string
	commandByte string
}
type MyConn struct {
	ID string
	net.Conn
	ignoreWrites bool
}

func readInput(conn net.Conn) (CommandInput, error){
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return CommandInput{}, err
	}

	command := [][]string {}
	input :=string(buf[:n])

	if(len(input) ==0){ 
		return CommandInput{
			commandStr: command,
			commandByte: input,
		}, nil
	}

	command, err = splitMultipleCommandString(input)

	if err != nil {
		return CommandInput{}, err
	}
	// switch(input[0]) {
	// case 43:
	// 	command = append(command,RESPSimpleString(input))
	// case 36:
	// 	command = RESPBulkString(input)
	// case 42:
	// 	command = RESPArray(input)
	// default:
	// 	fmt.Print("Unknown RESP enconfing")
	// }

	fmt.Println("what we printing here")
	fmt.Println(command)

	return CommandInput{
		commandStr: command,
		commandByte: input,
	}, nil

}


var replConn map[string]net.Conn = make(map[string]net.Conn)


func propagte (conn net.Conn, command []byte) {
	fmt.Println("propagte")
	if replConn == nil {
		return
	}
	fmt.Print(replConn)
	for _, v := range replConn {
		_, err := v.Write(command)

		if err != nil {
			fmt.Print("write error")
		}
	}
}

func (conn MyConn) Write(b []byte) (n int, err error) {
	if conn.ignoreWrites {
		return len(b), nil
	}
	return conn.Conn.Write(b)
}

var whitelistProp = map[Commands]bool{"SET": true}

func handleConenction(conn MyConn, serverCon Server) {
	defer func(conn MyConn) {
		fmt.Print("lets close it")
		// delete(replConn, conn.ID)
		conn.Close()
	}(conn)

	for {
		fmt.Println("read another")
		comamndInput, err := readInput(conn)
		if err != nil {
			fmt.Println("Error while reaidng", err)
			break;
		}
		
		for _,v := range comamndInput.commandStr {
		err = executeCommand(v,comamndInput.commandByte, conn, serverCon)
		if err != nil {
			fmt.Println("Error while reaidng", err)
			break;
		}
		}
		
	}
}

func executeCommand(commandI []string, bytes string, conn MyConn, serverCon Server) (err error){
	args := commandI

		if len(args) == 0 {
			fmt.Println("No argument passed")
			return err
		}

	command := Commands(strings.ToUpper(args[0]))
	if whitelistProp[command] && replConn != nil {
		propagte(conn, []byte(bytes))
	}

	switch(command) {
	case PING:
		err = conn.Ping()
	case ECHO:
		err = conn.Echo(args)
	case SET:	
		err = conn.Set(args)
	case GET:
		err = conn.Get(args)
	case INFO:
		err = conn.Info(args, serverCon)
	case REPLCONF:
		err = conn.ReplConf()
	case PSYNC: 
		err = conn.Psync(serverCon)
		replConn[conn.ID] = conn;
	default: {
		err = errors.New(fmt.Sprintf("invalid command received:%v", command))
		fmt.Println("invalid command received:", command)
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
	}
	
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}

	return err
}
