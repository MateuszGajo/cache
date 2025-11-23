package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"reflect"
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

func propagte(rawInput string, server *Server) {
	for _, v := range server.masterConfig.replicaConnections {
		v.bytesWrite += len(rawInput)
	}

	for _, v := range server.masterConfig.replicaConnections {
		_, err := v.Write([]byte(rawInput))

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
		respParsed, err := readInputNew(conn)

		if err != nil {
			fmt.Println("Error while reaidng", err)
			break
		}
	mainLoop:
		for _, respItem := range respParsed {
			// One array basically means one sets of command send, but we packets can be merge diffrently, so can be more than one array
			astArray, ok := respItem.astNode.(ASTArray)
			if !ok {
				log.Fatalf("error handling command, expected to get comamnd as ast array, got: %v", reflect.TypeOf(astArray))
				break
			}
			args := []string{}

			for _, arrayItem := range astArray.values {
				bulkNode, ok := arrayItem.(ASTBulkString)

				if !ok {
					log.Fatalf("all command in resp array should be bulk string, got: %v", reflect.TypeOf(bulkNode))
					break mainLoop
				}
				args = append(args, bulkNode.val)
			}

			if len(args) == 0 {
				log.Fatal("No argument passed")
				break
			}

			command := Commands(strings.ToUpper(args[0]))

			err := executeCommand(command, args[1:], conn, server)

			if err != nil {
				log.Fatal(err)
				break
			}

			whiteReplCommands := map[Commands]bool{SET: true}
			handshakeSyncCommands := map[Commands]bool{RDBFILE: true, FULLRESYNC: true}

			if whiteReplCommands[command] && len(server.masterConfig.replicaConnections) > 0 {
				propagte(respItem.rawInput, server)
			}

			if !handshakeSyncCommands[command] {
				server.replicaConfig.byteProcessed += len(respItem.rawInput)
			}
		}

	}
}

func executeCommand(command Commands, args []string, conn MyConn, server *Server) (err error) {

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
