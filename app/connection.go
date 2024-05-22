package main

import (
	"fmt"
	"net"
	"strings"
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
	FULLRESYNC Commands = "FULLRESYNC"
	RDBFILE Commands = "RDB-FILE"
	WAIT Commands = "WAIT"
)

type ReplicaConn struct {
	net.Conn
	bytesWrite int
	byteAck int
}


var whiteReplCommands = map[Commands]bool{SET: true}
var handshakeSyncCommands = map[Commands]bool{RDBFILE: true, FULLRESYNC: true}
var replConn map[string]*ReplicaConn = make(map[string]*ReplicaConn)
var byteParsed int = 0

func propagte (command CommandDetails) {
	fmt.Println("propagte")
	if replConn == nil {
		return
	}


	for _,v := range replConn {
		v.bytesWrite += command.length
		fmt.Println("co tutaj")
		fmt.Println(v.bytesWrite)
	}
	for _,v := range replConn {
		fmt.Println("co tutaj")
		fmt.Println(v.bytesWrite)
	}
	for _, v := range replConn {
		_, err := v.Write([]byte(command.raw))

		if err != nil {
			fmt.Print("write error")
		}
	}
}

func handleConenction(conn MyConn, serverCon Server) {
	defer func(conn MyConn) {
		fmt.Print("lets close it")
		delete(replConn, conn.ID)
		conn.Close()
	}(conn)

	for {
		comamndInput, err := readInput(conn)

		fmt.Println("read")
		fmt.Printf("%q",comamndInput)

		if err != nil {
			fmt.Println("Error while reaidng", err)
			break
		}

		for _, v := range comamndInput.commandStr {
			err = executeCommand(v, comamndInput.commandByte, conn, serverCon)

			if err != nil {
				fmt.Println("Error while reaidng", err)
				break
			}
		}

	}
}

func executeCommand(commandDetails  CommandDetails, bytes string, conn MyConn, serverCon Server) (err error){
	args := commandDetails.command

		if len(args) == 0 {
			fmt.Println("No argument passed")
			return err
		}

	command := Commands(strings.ToUpper(args[0]))


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
		err = conn.ReplConf(args)
	case PSYNC: 
		err = conn.Psync(serverCon)
		replConn[conn.ID] = &ReplicaConn{
			Conn: conn,
			bytesWrite: 0,
			byteAck: 0,
			};
	case FULLRESYNC:
		// empty for now
	case RDBFILE:
		// empty for now
	case WAIT:
		err = conn.Wait(args)
	default: {
		err = fmt.Errorf("invalid command received:%v", command)
		fmt.Println("invalid command received:", command)
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
	}
	if whiteReplCommands[command] && replConn != nil {
		propagte(commandDetails)
	}
	if(!handshakeSyncCommands[command]) {
		byteParsed += commandDetails.length
	}

	
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}

	return err
}