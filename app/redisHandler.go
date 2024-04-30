package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strconv"
)

var EMPTY_RDB_FILE_BASE64 string = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

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

func Psync(conn net.Conn,serverCon Server) string {
	_, err := conn.Write([]byte(BuildResponse("+FULLRESYNC " + serverCon.replicaId +" "+ strconv.Itoa(serverCon.replicaOffSet) +"\r\n")))

		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	emptyRdb, _ := base64.StdEncoding.DecodeString(EMPTY_RDB_FILE_BASE64)
	buff := []byte("$" + fmt.Sprint(len(emptyRdb)) + "\r\n")
	buff = append(buff, emptyRdb...)
	_, err  = conn.Write(buff)

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
	return ""
}