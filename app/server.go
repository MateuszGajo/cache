package main

import (
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




type CustomSetStore struct {
	Value      string
	ExpireAt  time.Time
}

type Replica struct {
	Port string
	Address string
}

type Server struct {
	role string
	replicaId string
	replicaOffSet int
}

type MyConn struct {
	ID string
	net.Conn
	ignoreWrites bool
}

var port int
var replica Replica
var CLRF string

func init(){
	flag.IntVar(&port, "port", 6379, "port to listen to")
	replicaof := flag.String("replicaof", "", "Replica of master address")

	flag.Parse()
	if *replicaof != "" {
		parts := strings.Split(*replicaof, " ")
		if len(parts) == 2 {
			replica.Address = parts[0]
			replica.Port = parts[1]
		}
	}

	os := runtime.GOOS

    if os == "windows" {
        CLRF = "\\r\\n" // fix this
    } else {
        CLRF = "\r\n"
    }
}


func handShake(){
	conn, err := net.Dial("tcp", replica.Address + ":" + replica.Port)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", replica.Address, replica.Port)
	}

	_, err = conn.Write([]byte(BuildRESPArray([]string{"ping"})))

	if err != nil {
		fmt.Print("error while pinging master replica", err)
		conn.Close()
		return
	}

	inputComm, err := readInput(conn)
	if err != nil {
		fmt.Print("error while pinging master replica", err)
		conn.Close()
		return
	}

	args := inputComm.commandStr[0].command
	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte(BuildRESPArray([]string{"REPLCONF","listening-port",strconv.Itoa((port))})))

	inputComm, err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf listening-port master replica")
		conn.Close()
		return
	}

	args = inputComm.commandStr[0].command
	if args[0] != "OK" {
		fmt.Printf("Response its invalid, expected ok we got:%v", args[0])
		os.Exit(1)
	}

	conn.Write([]byte(BuildRESPArray([]string {"REPLCONF","capa","psync2"})))

	inputComm,err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf capa  master replica")
		conn.Close()
		return
	}

	args = inputComm.commandStr[0].command
	if args[0] != "OK" {
		fmt.Printf("Response its invalid, expected ok we got:%v", args[0])
		os.Exit(1)
	}

	conn.Write([]byte(BuildRESPArray([]string{"PSYNC","?","-1"})))

	go handleConenction(MyConn{Conn: conn, ignoreWrites: false, ID: strconv.Itoa(rand.IntN(100))}, Server{}) 
}

func main() {
	serverCon := Server{
		role: "master",
	}
	

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	if err != nil {
		fmt.Printf("Failed to run on port %d", port)
		os.Exit(1)
	}

	defer func(){
		fmt.Print("close")
		ln.Close()
	}()

	if(replica != Replica{}) {
		serverCon.role = "slave"
		handShake()
	} else {
		serverCon.replicaOffSet = 0
		serverCon.replicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	}


	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			os.Exit(1)
		}

		go handleConenction(MyConn{Conn: conn, ignoreWrites: false, ID: strconv.Itoa(rand.IntN(100)) }, serverCon)
	}
	
}














