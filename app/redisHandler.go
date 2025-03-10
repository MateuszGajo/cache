package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var EMPTY_RDB_FILE_BASE64 string = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

// find a way to not respond to master command
// write integration test
func (conn MyConn) Ping() (err error) {

	//let thing fo better solution
	connectionFromMaster := strings.Contains(conn.RemoteAddr().String(), "6379")
	if !connectionFromMaster {
		_, err = conn.Write([]byte(BuildSimpleString("PONG")))
		return err
	}
	return nil
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
	if len(args) >= 4 {
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
		opResult = handleSet(key, value, &timeMs, "string")
	default:
		opResult = handleSet(key, value, nil, "string")
	}

	connectionFromMaster := strings.Contains(conn.RemoteAddr().String(), "6379")

	if opResult && !connectionFromMaster {
		conn.Write([]byte(BuildSimpleString("OK")))
	} else if !opResult {
		return errors.New("Problem setting variable")
	}

	return nil
}

func (conn MyConn) Get(args []string) (err error) {
	key := args[1]
	var result string
	value := handleGet(key)

	if (value != CustomSetStore{}) {
		result = BuildBulkString(value.Value)
	} else {
		result = BuildNullBulkString()
	}

	_, err = conn.Write([]byte(result))
	return err

}

func (conn MyConn) Info(args []string, serverCon Server) (err error) {
	result := BuildBulkString(fmt.Sprintf("role:%vmaster_replid:%vmaster_repl_offset:%v", serverCon.role, serverCon.replicaId, strconv.Itoa(serverCon.replicaOffSet)))
	_, err = conn.Write([]byte(result))
	return err
}

type ReplConfCommand string

const (
	LISTENING_PORT ReplConfCommand = "LISTENING-PORT"
	CAPA           ReplConfCommand = "CAPA"
	GETACK         ReplConfCommand = "GETACK"
	ACK            ReplConfCommand = "ACK"
)

func (conn MyConn) replConfConfirm() (err error) {
	fmt.Println("repl cofnirm")
	result := BuildSimpleString("OK")
	fmt.Println("simple string")
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn) replConfGetAct() (err error) {
	result := BuildRESPArray([]string{"REPLCONF", "ACK", strconv.Itoa(byteParsed)})
	_, err = conn.Write([]byte(result))
	return err
}

func (conn MyConn) replConfAct(args []string) (err error) {
	bytes := args[2]
	for _, rc := range replConn {
		if rc.RemoteAddr().String() == conn.RemoteAddr().String() {
			fmt.Println("act bytes for")
			fmt.Print(rc.RemoteAddr().String())
			fmt.Fprintln(os.Stdout, []any{"\n, bytes %v", bytes}...)
			conv, _ := strconv.Atoi(bytes)
			rc.byteAck = conv
		}
	}
	return nil
}

func (conn MyConn) ReplConf(args []string) (err error) {
	command := ReplConfCommand(strings.ToUpper(args[1]))

	fmt.Println("args", args)
	fmt.Println("command", command)
	fmt.Println(command == LISTENING_PORT)

	switch command {
	case LISTENING_PORT, CAPA:
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
	resync := fmt.Sprintf("+FULLRESYNC %v %v%v", serverCon.replicaId, strconv.Itoa(serverCon.replicaOffSet), CLRF)
	fmt.Println("sedning resync")
	_, err = conn.Write([]byte(resync))
	if err != nil {
		return err
	}

	emptyRdb, _ := base64.StdEncoding.DecodeString(EMPTY_RDB_FILE_BASE64)
	rdbBuffer := []byte(fmt.Sprintf("$%v%v", fmt.Sprint(len(emptyRdb)), CLRF))
	rdbBuffer = append(rdbBuffer, emptyRdb...)

	_, err = conn.Write([]byte(rdbBuffer))
	if err != nil {
		fmt.Print("problem with file")
	}
	return err
}

func (conn MyConn) Type(args []string) (err error) {

	key := args[1]
	var result string
	value := handleGet(key)

	if (value != CustomSetStore{}) {
		result = BuildSimpleString(value.Type)
	} else {
		result = BuildSimpleString("none")
	}

	_, err = conn.Write([]byte(result))
	return err
}

func Deserialize(input string) [][]string {
	if(input == "") {
		return [][]string{}
	}

	fmt.Println("input")
	fmt.Println(input)

	entries := strings.Split(input, ",")

	result := make([][]string, 0, len(entries)*2)
	for i := 0; i < len(entries); i++ {
		split := strings.Split(entries[i], ":")
		entry := []string {split[0],split[1]}
		result = append(result, entry)
	}

	return result
}

func GetLastStreamEntry(stream [][]string) (result [][]string) {
	for i:= len(stream) -1; i >= 0; i-- {
		if(stream[i][0] == "id") {
			result = stream[i:]
			return result
		}
	}

	return result
}


func VerifyStreamIntegrity(stream [][]string, newEntryId string) (bool,string) {
	currentEntryIdSplited := strings.Split(newEntryId, "-")
	currentEntryKeyTime := currentEntryIdSplited[0]
	currnetEntrySequenceNumber := currentEntryIdSplited[1]

	if(currentEntryKeyTime == "0" && currnetEntrySequenceNumber == "0") {
		return true,"The ID specified in XADD must be greater than 0-0"
	}

	lastEntry := GetLastStreamEntry(stream)
	lastEntryKeyTime := "0"
	lastEntrySequenceNumber := "0"
	fmt.Println("stream")
	fmt.Println(stream)
	fmt.Println("last entry")
	fmt.Println(lastEntry)
	if len(lastEntry) >0 {
		lastEntryId := strings.Split(lastEntry[0][1], "-")
		lastEntryKeyTime = lastEntryId[0]
		lastEntrySequenceNumber = lastEntryId[1]
	}

	if(lastEntryKeyTime == currentEntryKeyTime) {
		if(currnetEntrySequenceNumber >lastEntrySequenceNumber) {
			return false, ""
		} else {
			return true, "The ID specified in XADD is equal or smaller than the target stream top item"
		}
	} else if (currentEntryKeyTime > lastEntryKeyTime) {
		return false, ""
	}

	return true, "The ID specified in XADD is equal or smaller than the target stream top item"
}

func (conn MyConn) Xadd(args []string) (err error) {
	streamKey := args[1]
	entryKey := args[2]

	previousEntry := handleGet(streamKey)

	isErr, errMsg := VerifyStreamIntegrity(Deserialize(previousEntry.Value),entryKey);
	if isErr {
		conn.Write([]byte(BuildSimpleError("ERR", errMsg)))
		return;
	}

	content := []string{fmt.Sprintf("id:%v", entryKey)}
	if len(args)%2 == 0 {
		fmt.Print("wrong number of arguments")
		return
	}
	for i := 3; i < len(args)-1; i++ {
		content = append(content, fmt.Sprintf("%v:%v", args[i], args[i+1]))
	}
	serializedString := strings.Join(content, ",")

	isErr = handleSet(streamKey, serializedString, nil, "stream")

	if !isErr {
		fmt.Print("Problem setting value")
		return
	}

	conn.Write([]byte(BuildBulkString(entryKey)))

	return err
}

func (conn MyConn) XaddHandler(args []string) (bool, string) {
	streamKey := args[1]
	entryKey := args[2]

	previousEntry := handleGet(streamKey)

	ok, errMsg := VerifyStreamIntegrity(Deserialize(previousEntry.Value),entryKey);
	if !ok {
		// conn.Write([]byte())
		return false, BuildSimpleError("ERR", errMsg);
	}

	content := []string{fmt.Sprintf("id:%v", entryKey)}
	if len(args)%2 == 0 {
		fmt.Print("wrong number of arguments")
		return true, "wrong number of arguments"
	}
	for i := 3; i < len(args)-1; i++ {
		content = append(content, fmt.Sprintf("%v,%v", args[i], args[i+1]))
	}
	serializedString := strings.Join(content, ",")

	ok = handleSet(streamKey, serializedString, nil, "stream")

	if !ok {
		fmt.Print("Problem setting value")
		return true, "PRoblem setting value"
	}

	return false, BuildBulkString(entryKey)
}

func (conn MyConn) Wait(args []string) (err error) {

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

	getAckMsg := BuildRESPArray([]string{"REPLCONF", "GETACK", "*"})
	getAckMsgLen := len(getAckMsg)
	for _, v := range replConn {
		go func(v *ReplicaConn) {
			if v.bytesWrite > 0 {
				v.Write([]byte(getAckMsg))
				v.bytesWrite += getAckMsgLen
			}
		}(v)
	}

	timer := time.NewTimer(time.Duration(timeoutInMs) * time.Millisecond)
	ticker := time.NewTicker(100 * time.Millisecond)

	defer ticker.Stop()

	totalAcked := 0
	for {
		select {
		case <-timer.C:
			fmt.Println("timer end1")
			conn.Write([]byte(BuildRespInt(totalAcked)))
			return
		case <-ticker.C:
			totalAcked = 0
			for _, v := range replConn {
				bytesWritten := max(0, v.bytesWrite-getAckMsgLen)
				if bytesWritten == v.byteAck {
					totalAcked++
				}
			}
			if totalAcked >= acksRequired {
				conn.Write([]byte(BuildRespInt(totalAcked)))
				return
			}

		}
	}
}
