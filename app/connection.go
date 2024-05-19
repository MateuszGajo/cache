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


var whiteReplCommands = map[Commands]bool{SET: true}
var handshakeSyncCommands = map[Commands]bool{RDBFILE: true, FULLRESYNC: true}
var replConn map[string]net.Conn = make(map[string]net.Conn)
var byteParsed int = 0

func propagte (command []byte) {
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

func handleConenction(conn MyConn, serverCon Server) {
	defer func(conn MyConn) {
		fmt.Print("lets close it")
		delete(replConn, conn.ID)
		conn.Close()
	}(conn)

	for {
		comamndInput, err := readInput(conn)

		fmt.Println("read")
		fmt.Println(comamndInput)

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

	if whiteReplCommands[command] && replConn != nil {
		propagte([]byte(bytes))
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
		err = conn.ReplConf(args)
	case PSYNC: 
		err = conn.Psync(serverCon)
		replConn[conn.ID] = conn;
	case FULLRESYNC:
		// empty for now
	case RDBFILE:
		// empty for now
	case WAIT:
		err = conn.Wait()
	default: {
		err = fmt.Errorf("invalid command received:%v", command)
		fmt.Println("invalid command received:", command)
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
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