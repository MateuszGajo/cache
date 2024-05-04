package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
)

var EMPTY_RDB_FILE_BASE64 string = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

func (conn MyConn) Ping() (err error) {
	_, err = conn.Write([]byte(BuildSimpleString("PONG")))
	return err
}

func (conn MyConn) Echo(args []string) (err error) {
	msg := args[1]
	_, err = conn.Write([]byte(BuildBulkString(msg)))
	return err
}

func (conn MyConn) Set(args []string) (err error) {
	key := args[1]
	value := args[2]
	param, paramValue := "", ""
	var opResult bool
	if(len(args) >=4) {
		param = args[3]
		paramValue = args[4]
	}

	switch param {
	case "px":
		timeMs, err := strconv.Atoi(paramValue)
		if err != nil {
			return errors.New("Invalid time")
		}
		opResult = handleSet(key, value, &timeMs)
	default:
		opResult = handleSet(key, value, nil)
	}

	if opResult {
		conn.Write([]byte(BuildSimpleString("OK")))
	} else {
		return errors.New("Problem setting variable")
	}

	return nil;
}

func (conn MyConn) Get(args [] string) (err error) {
	key := args[1]
	var value, result string
	value, err = handleGet(key)

	if err != nil {
		return err
	}

	if(value != "") {
		result = BuildBulkString(value)
	} else {
		result = BuildNullBulkString()
	}

	_, err = conn.Write([]byte(result))
	return err

}

func (conn MyConn) Info(args []string, serverCon Server) (err error) {
	result := BuildBulkString(fmt.Sprintf("role:%vmaster_replid:%vmaster_repl_offset:%v",serverCon.role, serverCon.replicaId, strconv.Itoa(serverCon.replicaOffSet)))
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn)  ReplConf() (err error) {
	result :=  BuildSimpleString("OK")
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn) Psync(serverCon Server) (err error) {
	resync := fmt.Sprintf("+FULLRESYNC %v %v%v",serverCon.replicaId,strconv.Itoa(serverCon.replicaOffSet),CLRF)
	_, err = conn.Write([]byte(resync))
	if err != nil {
		return err
	}

	emptyRdb, _ := base64.StdEncoding.DecodeString(EMPTY_RDB_FILE_BASE64)
	rdbBuffer := []byte(fmt.Sprintf("$%v%v",fmt.Sprint(len(emptyRdb)), CLRF))
	rdbBuffer = append(rdbBuffer, emptyRdb...)

	_, err = conn.Write([]byte(rdbBuffer))
	return err
}