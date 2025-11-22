package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// this RESPDATA, RESPIntData are going to be removed
type RESPData interface{}

type RESPIntData int

type RESPRecord struct {
	data   RESPData
	length int
	raw    string
}

var tempRead string = ""

type RESPParsed struct {
	records []RESPRecord
	raw     string
}

type RESPParsedNew struct {
	nodes []ASTNode
	raw   string
}

// would be good to create a struct reader and incilized with parser

func readInputNew(conn net.Conn) (RESPParsedNew, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return RESPParsedNew{}, err
	}

	nodes := []ASTNode{}
	input := string(buf[:n])

	if len(input) == 0 {
		return RESPParsedNew{
			nodes: nodes,
			raw:   input,
		}, nil
	}

	parser := NewParser()

	parsed := parser.parseStream(input)

	if parsed.err != nil {
		return RESPParsedNew{}, parsed.err
	}

	return RESPParsedNew{
		nodes: parsed.AST,
		raw:   input,
	}, nil

}

func readInput(conn net.Conn) (RESPParsed, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return RESPParsed{}, err
	}

	command := []RESPRecord{}
	input := string(buf[:n])

	if len(input) == 0 {
		return RESPParsed{
			records: command,
			raw:     input,
		}, nil
	}

	fmt.Printf("what we got here? :%q\n", tempRead+input)

	command, tempRead, err = splitMultipleCommandString(tempRead + input)
	fmt.Printf("resp command: %+v \n", command)

	if err != nil {
		return RESPParsed{}, err
	}

	return RESPParsed{
		records: command,
		raw:     input,
	}, nil

}

func splitMultipleCommandString(input string) (res []RESPRecord, tempInp string, err error) {
	tempInp = input
	endIndex := 0
	var record RESPRecord

loop:
	for len(tempInp) > 0 {
		indicator := tempInp[0]

		switch indicator {
		case '+':
			endIndex, record = handleSimpleString(tempInp)
		case '$':
			endIndex, record = handleBulkString(tempInp)
		case '*':
			endIndex, record = handleRespArray(tempInp)
		case '-':
			endIndex, record = handleSimpleError(tempInp)
		case ':':
			endIndex, record = handleInteger(tempInp)
		default:
			fmt.Println("unsupported break: ", indicator)
			break loop
		}
		if endIndex == -1 {
			break loop
		}
		res = append(res, record)
		tempInp = tempInp[endIndex:]
	}

	return res, tempInp, nil
}

func handleInteger(input string) (endIndex int, command RESPRecord) {
	index := strings.Index(input, CLRF)
	// i dont like adding index == -1 in every handler
	if index == -1 {
		endIndex = -1
		return endIndex, command
	}
	endIndex = index + len(CLRF)

	integer, err := strconv.Atoi(input[1:index])

	if err != nil {
		panic(err)
	}

	command = RESPRecord{
		data:   RESPIntData(integer),
		length: len(input[:endIndex]),
		raw:    input[:endIndex],
	}

	return endIndex, command
}

func handleSimpleError(input string) (endIndex int, command RESPRecord) {
	index := strings.Index(input, CLRF)
	if index == -1 {
		endIndex = -1
		return endIndex, command
	}
	endIndex = index + len(CLRF)

	command = RESPRecord{
		data:   []string{input[:index]},
		length: len(input[:endIndex]),
		raw:    input[:endIndex],
	}

	return endIndex, command
}

func handleSimpleString(input string) (endIndex int, command RESPRecord) {
	index := strings.Index(input, CLRF)
	if index == -1 {
		fmt.Println("inside?")
		endIndex = -1
		return endIndex, command
	}

	endIndex = index + len(CLRF)
	if strings.Contains(input, "FULLRESYNC") {
		command = RESPRecord{
			data:   strings.Split(input[1:index], " "),
			length: len(input[:endIndex]),
			raw:    input[endIndex:],
		}
	} else {
		command = RESPRecord{
			data:   input[1:index],
			length: len(input[:endIndex]),
			raw:    input[:endIndex],
		}
	}

	return endIndex, command
}

func handleBulkString(input string) (endIndex int, command RESPRecord) {
	index := strings.Index(input, CLRF)
	if index == -1 {
		endIndex = -1
		return endIndex, command
	}

	commLen, _ := strconv.Atoi(input[1:index])

	delimiterLen := len(CLRF)
	endIndex = commLen + index + delimiterLen

	if len(input) >= endIndex+delimiterLen && input[endIndex] == 13 && input[endIndex+1] == 10 {

		command = RESPRecord{
			data:   []string{input[index+delimiterLen : endIndex]},
			length: len(input[:endIndex]),
			raw:    input[:endIndex],
		}
		endIndex += 2
	} else {
		command = RESPRecord{
			data:   []string{"rdb-file", input[index+delimiterLen : endIndex]},
			length: len(input[:endIndex]),
			raw:    input[:endIndex],
		}
	}

	return endIndex, command
}

func handleRespArray(input string) (endIndex int, command RESPRecord) {
	index := strings.Index(input, CLRF)
	if index == -1 {
		endIndex = -1
		return endIndex, command
	}

	commLen, _ := strconv.Atoi(input[1:index])
	endIndex = findOccurance(input, "\r", commLen*2)
	if endIndex == -1 {
		endIndex = -1
		return endIndex, command
	}

	stringToDo := input[:endIndex]
	parts := strings.Split(stringToDo, CLRF)
	temp := []string{}
	for i := 2; i < len(parts); i += 2 {
		temp = append(temp, parts[i])
	}

	command = RESPRecord{
		data:   temp,
		length: len(input[:endIndex+2]),
		raw:    input[:endIndex+2],
	}

	return endIndex + 2, command
}

func findOccurance(input string, delimiter string, occNum int) int {
	re := regexp.MustCompile(delimiter)

	indexes := re.FindAllStringIndex(input, -1)

	if len(indexes) < occNum {
		return -1
	}

	return indexes[occNum][0]

}
