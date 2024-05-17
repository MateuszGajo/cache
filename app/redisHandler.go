package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var EMPTY_RDB_FILE_BASE64 string = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

func (conn MyConn) Ping() (err error) {
	_, err = conn.Write([]byte(BuildSimpleString("PONG")))
	return err
}

func (conn MyConn)   Echo(args []string) (err error) {
	msg := args[1]
	_, err = conn.Write([]byte(BuildBulkString(msg)))
	return err
}

func (conn MyConn)  Set(args []string) (err error) {
	key := args[1]
	value := args[2]
	param, paramValue := "", ""
	var opResult bool
	if(len(args) >=4) {
		param = args[3]
		paramValue = args[4]
	}
	fmt.Print("value")
	fmt.Print(value)

	switch param {
	case "px":
		timeMs, err := strconv.Atoi(paramValue)
		if err != nil {
			return errors.New("Invalid time")
		}
		fmt.Print("set 2111")
		opResult = handleSet(key, value, &timeMs)
	default:
		fmt.Print("set 32")
		opResult = handleSet(key, value, nil)
	}

	connectionFromMaster := strings.Contains(conn.RemoteAddr().String(), "6379")

	if opResult && !connectionFromMaster {
		conn.Write([]byte(BuildSimpleString("OK")))
	} else if !opResult {
		return errors.New("Problem setting variable")
	}

	return nil;
}

func (conn MyConn)  Get(args [] string) (err error) {
	key := args[1]
	var value, result string
	value = handleGet(key)

	fmt.Println("what is the value?")
	fmt.Println(value)
	if(value != "") {
		fmt.Println("not empty value")
		result = BuildBulkString(value)
	} else {
		fmt.Println("empty value")
		result = BuildNullBulkString()
	}

	fmt.Println("what are we sending?")
	fmt.Println(result)

	_, err = conn.Write([]byte(result))
	return err

}

func (conn MyConn)  Info(args []string, serverCon Server) (err error) {
	result := BuildBulkString(fmt.Sprintf("role:%vmaster_replid:%vmaster_repl_offset:%v",serverCon.role, serverCon.replicaId, strconv.Itoa(serverCon.replicaOffSet)))
	_, err = conn.Write([]byte(result))
	return err
}

type ReplConfCommand string
const (
	LISTENING_PORT ReplConfCommand = "LISTENING-PORT"
	CAPA ReplConfCommand = "CAPA"
	GETACK ReplConfCommand = "GETACT"
)



func (conn MyConn)  replConfConfirm() (err error) {
	result :=  BuildSimpleString("OK")
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn)  replConfAct() (err error) {
	result :=  BuildRESPArray([]string{"REPLCONF", "ACK", "0"})
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn)  ReplConf(args []string) (err error) {
	command := ReplConfCommand(strings.ToUpper(args[0]))

	switch (command) {
	case LISTENING_PORT:
	case CAPA:
		return conn.replConfConfirm()
	case GETACK:
		return conn.replConfAct()
	default:
		fmt.Printf("could find command: %v with args: %v", command, args)
	}
	result :=  BuildSimpleString("OK")
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn) Psync(serverCon Server) (err error) {
	resync := fmt.Sprintf("+FULLRESYNC %v %v%v",serverCon.replicaId,strconv.Itoa(serverCon.replicaOffSet),CLRF)
	fmt.Println("sedning resync")
	_, err = conn.Write([]byte(resync))
	if err != nil {
		return err
	}
	

	emptyRdb, _ := base64.StdEncoding.DecodeString(EMPTY_RDB_FILE_BASE64)
	rdbBuffer := []byte(fmt.Sprintf("$%v%v",fmt.Sprint(len(emptyRdb)), CLRF))
	rdbBuffer = append(rdbBuffer, emptyRdb...)

	_, err = conn.Write([]byte(rdbBuffer))
	if err !=nil {
		fmt.Print("problem with file")
	}
	return err
}