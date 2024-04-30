package main

import (
	"fmt"
	"os"
	"strconv"
)

func Ping() string {
	return "+PONG\r\n"
}

func Echo(args []string) string {
	return BuildResponse(args[1])
}

func Set(args []string) string {
	switch len(args) {
	case 5: //THAT MAGIC VALUE FOR
		// command := args[3]
		timeMs, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Print("invalid time")
			os.Exit(1)
		}
		return handleSet(args[1], args[2], &timeMs)
	case 3:
		return handleSet(args[1], args[2], nil)
	}

	return ""
}

func Get(args [] string) string {
	return handleGet(args[1])
}

func Info(args []string, serverCon Server) string {
	return BuildResponse("role:"+ serverCon.role +"master_replid:" + serverCon.replicaId +"master_repl_offset:" + strconv.Itoa(serverCon.replicaOffSet))
}

func ReplConf() string {
	return "+OK\r\n"
}

func Psync(serverCon Server) string {
	return BuildResponse("+FULLRESYNC " + serverCon.replicaId +" "+ strconv.Itoa(serverCon.replicaOffSet) +"\r\n")
}