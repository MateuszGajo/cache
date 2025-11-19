package main

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"strings"
)

type Commands string

const (
	PING       Commands = "PING"
	ECHO       Commands = "ECHO"
	SET        Commands = "SET"
	GET        Commands = "GET"
	INFO       Commands = "INFO"
	REPLCONF   Commands = "REPLCONF"
	PSYNC      Commands = "PSYNC"
	FULLRESYNC Commands = "FULLRESYNC"
	RDBFILE    Commands = "RDB-FILE"
	WAIT       Commands = "WAIT"
	TYPE       Commands = "TYPE"
	RPUSH      Commands = "RPUSH"
	XADD       Commands = "XADD"
	XRANGE     Commands = "XRANGE"
	XREAD      Commands = "XREAD"
	CONFIG     Commands = "CONFIG"
	KEYS       Commands = "KEYS"
)

func propagte(command RESPRecord, server *Server) {
	for _, v := range server.masterConfig.replicaConnections {
		v.bytesWrite += command.length
	}

	for _, v := range server.masterConfig.replicaConnections {
		_, err := v.Write([]byte(command.raw))

		if err != nil {
			fmt.Print("write error")
		}
	}
}

func handleConenction(conn MyConn, server *Server) {
	defer func(conn MyConn) {
		fmt.Print("lets close it")
		delete(server.masterConfig.replicaConnections, conn.ID)
		conn.Close()
	}(conn)

	for {
		comamndInput, err := readInput(conn)

		if err != nil {
			fmt.Println("Error while reaidng", err)
			break
		}

		for _, v := range comamndInput.records {
			err = executeCommand(v, comamndInput.raw, conn, server)

			if err != nil {
				fmt.Println("Error while reaidng", err)
				break
			}
		}

	}
}

func executeCommand(commandDetails RESPRecord, bytes string, conn MyConn, server *Server) (err error) {
	args := commandDetails.data.([]string)

	if len(args) == 0 {
		fmt.Println("No argument passed")
		return err
	}

	command := Commands(strings.ToUpper(args[0]))

	switch command {
	case PING:
		err = conn.Ping()
	case ECHO:
		err = conn.Echo(args)
	case SET:
		err = conn.Set(args)
	case GET:
		err = conn.Get(args, server)
	case INFO:
		err = conn.Info(args, server)
	case REPLCONF:
		err = conn.ReplConf(args, server)
	case TYPE:
		err = conn.Type(args)
	case RPUSH:
		err = conn.RPush(args, server)
	case XADD:
		err = conn.Xadd(args)
	case XRANGE:
		err = conn.Xrange(args)
	case XREAD:
		err = conn.Xread(args)
	case KEYS:
		err = conn.Keys(args, server)
	case CONFIG:
		err = conn.Config(args, server)
	case PSYNC:
		err = conn.Psync(server)
		server.masterConfig.replicaConnections[conn.ID] = &ReplicaConn{
			Conn:       conn,
			bytesWrite: 0,
			byteAck:    0,
		}
	case FULLRESYNC:
		// empty for now
	case RDBFILE:
		// empty for now
	case WAIT:
		err = conn.Wait(args, server)
	default:
		{
			err = fmt.Errorf("invalid command received:%v", command)
			fmt.Println("invalid command received:", command)
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
	whiteReplCommands := map[Commands]bool{SET: true}
	handshakeSyncCommands := map[Commands]bool{RDBFILE: true, FULLRESYNC: true}

	if whiteReplCommands[command] && len(server.masterConfig.replicaConnections) > 0 {
		propagte(commandDetails, server)
	}

	if !handshakeSyncCommands[command] {
		server.replicaConfig.byteProcessed += commandDetails.length
	}

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return err
	}

	return err
}

func handShake(server *Server) {
	fmt.Println("hello")
	conn, err := net.Dial("tcp", server.replicaConfig.masterAddress+":"+server.replicaConfig.masterPort)

	if err != nil {
		fmt.Printf("cannot connect to %v:%v", server.replicaConfig.masterAddress, server.replicaConfig.masterPort)
	}

	_, err = conn.Write([]byte(BuildPrimitiveRESPArray([]string{"ping"})))

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

	args := inputComm.records[0].data.([]string)
	if args[0] != "PONG" {
		fmt.Print("Response its invalid")
		os.Exit(1)
	}

	conn.Write([]byte(BuildPrimitiveRESPArray([]string{"REPLCONF", "listening-port", server.port})))

	inputComm, err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf listening-port master replica")
		conn.Close()
		return
	}

	args = inputComm.records[0].data.([]string)
	if args[0] != "OK" {
		fmt.Printf("Response its invalid, expected ok we got:%v", args[0])
		os.Exit(1)
	}

	conn.Write([]byte(BuildPrimitiveRESPArray([]string{"REPLCONF", "capa", "psync2"})))

	inputComm, err = readInput(conn)
	if err != nil {
		fmt.Print("error while replConf capa  master replica")
		conn.Close()
		return
	}

	args = inputComm.records[0].data.([]string)
	if args[0] != "OK" {
		fmt.Printf("Response its invalid, expected ok we got:%v", args[0])
		os.Exit(1)
	}

	conn.Write([]byte(BuildPrimitiveRESPArray([]string{"PSYNC", "?", "-1"})))

	go handleConenction(MyConn{Conn: conn, ignoreWrites: false, ID: strconv.Itoa(rand.IntN(100))}, &Server{})
}
