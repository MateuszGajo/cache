package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var EMPTY_RDB_FILE_BASE64 string = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
// find a way to not respond to master command
// write integration test
func (conn MyConn) Ping() (err error) {

	//let thing o better solution
	connectionFromMaster := strings.Contains(conn.RemoteAddr().String(), "6379")
	if !connectionFromMaster {
		_, err = conn.Write([]byte(BuildSimpleString("PONG")))
		return err
	}
	return nil
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
	GETACK ReplConfCommand = "GETACK"
	ACK ReplConfCommand = "ACK"
)


func (conn MyConn)  replConfConfirm() (err error) {
	fmt.Println("repl cofnirm")
	result :=  BuildSimpleString("OK")
	fmt.Println("simple string")
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn)  replConfGetAct() (err error) {
	fmt.Println("This should be first")
	result :=  BuildRESPArray([]string{"REPLCONF", "ACK",  strconv.Itoa(byteParsed)})
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn)  replConfAct(args []string) (err error) {
	fmt.Println("This should be first")
	fmt.Println("args")
	fmt.Println(args)
	bytes := args[2]
	for _, rc := range replConn {
		if rc.RemoteAddr().String() == conn.RemoteAddr().String() {
			conv, _ := strconv.Atoi(bytes)
			rc.byteAck = conv
		}
	}
	// result :=  BuildRESPArray([]string{"REPLCONF", "ACK",  strconv.Itoa(byteParsed)})
	// _, err = conn.Write([]byte(result))
	return nil
}

func (conn MyConn)  ReplConf(args []string) (err error) {
	command := ReplConfCommand(strings.ToUpper(args[1]))

	fmt.Println("args", args)
	fmt.Println("command", command)
	fmt.Println(command == LISTENING_PORT)

	switch (command) {
	case LISTENING_PORT,CAPA:
		return conn.replConfConfirm()
	case GETACK:
		return conn.replConfGetAct()
	case ACK:
		return conn.replConfAct(args)
	default:
		fmt.Printf("could find command: %v with args: %v", command, args)
	}
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

func (conn MyConn) Wait(args []string) (err error) {
	// build wait
	conn.Write([]byte(BuildRespInt(len(replConn))))

	acksRequired, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	timeoutInMs, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}


	for _, v := range replConn {
		go func(v net.Conn){
			v.Write([]byte(BuildRESPArray([]string{"REPLCONF", "GETACK", "*"})))
		}(v)
	}


	timer := time.NewTimer(time.Duration(timeoutInMs) * time.Millisecond)
	ticker := time.NewTicker(100 * time.Millisecond)

	defer ticker.Stop()

	totalAcked := 0
	for {
		select {
		case <- timer.C:
			fmt.Println("timer end1")
			conn.Write([]byte(BuildRespInt(acksRequired)))
			return
		case <- ticker.C:
			totalAcked = 0
			for _, v := range replConn{
				if(v.bytesWrite <= v.byteAck){
					fmt.Println("total act +1")
					fmt.Printf("Bytes write, %v \n", v.bytesWrite)
					fmt.Printf("Bytes ack, %v \n", v.byteAck)
					totalAcked++
				}
			}
			if(totalAcked >= acksRequired){
				fmt.Println("all act good")
				fmt.Println(acksRequired)
				conn.Write([]byte(BuildRespInt(acksRequired)))
				return
			}
			
		}
	}
		
	return err
}