package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type CommandDetails struct {
	command []string
	length int
	raw string
}

var tempRead string = ""

type CommandInput struct{
	commandStr []CommandDetails
	commandByte string
}

func readInput(conn net.Conn) (CommandInput, error){
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Print("what we read", string(buf[:n]))
		return CommandInput{}, err
	}

	command := []CommandDetails {}
	input :=string(buf[:n])

	if(len(input) ==0){ 
		return CommandInput{
			commandStr: command,
			commandByte: input,
		}, nil
	}

	command, tempRead, err = splitMultipleCommandString(tempRead + input)

	if err != nil {
		return CommandInput{}, err
	}

	return CommandInput{
		commandStr: command,
		commandByte: input,
	}, nil

}


func splitMultipleCommandString(input string) (res []CommandDetails, tempInp string, err error) {
	tempInp = input
	endIndex := 0
	var command CommandDetails
	

	loop: for (len(tempInp) > 0){
		indicator := tempInp[0]

	switch (indicator) {
	case 43:
		endIndex, command = handleSimpleString(tempInp)
			
	case 36:
		endIndex, command = handleBulkString(tempInp)

	case 42:
		endIndex, command = handleRespArray(tempInp)
	default:
		fmt.Println("unsupported break")
		break loop;
	}
	if(endIndex == -1) {
		break loop;
	}
	res= append(res, command)
	tempInp = tempInp[endIndex:]
}


	return res,tempInp, nil
}

func handleSimpleString(input string) (endIndex int, command CommandDetails) {
	index := strings.Index(input, CLRF)
	if(index == -1) {
		endIndex = -1
		return endIndex, command
	}

	endIndex = index+ len(CLRF)
	if strings.Contains(input, "FULLRESYNC") {
		command = CommandDetails{
			command: strings.Split(input[1:index], " "),
			length: len(input[:endIndex]),
			raw: input[endIndex:],
		}
	}else {
		command = CommandDetails{
			command: []string{input[1:index]},
			length: len(input[:endIndex]),
			raw: input[:endIndex],
		}
	}
	
	return endIndex, command
}

func handleBulkString(input string) (endIndex int, command CommandDetails) {
		index := strings.Index(input, CLRF)
		if(index == -1) {
			endIndex = -1
			return endIndex, command
		}

		commLen, _ := strconv.Atoi(input[1:index])

		delimiterLen := len(CLRF)
		endIndex = commLen + index + delimiterLen

		if(len(input) >= endIndex + delimiterLen && input[endIndex +1] == 13 && input[endIndex +2] == 10) {
			//normal command
			endIndex +=2
			
			command = CommandDetails{
				command: []string{input[index+delimiterLen:endIndex]},
				length: len(input[:endIndex]),
				raw: input[:endIndex],
			}
		} else {
			//rdb file			

			command = CommandDetails{
				command: []string{"rdb-file", input[index+ delimiterLen:endIndex]},
				length: len(input[:endIndex]),
				raw: input[:endIndex],
			}
		}

		return endIndex, command
}

func handleRespArray(input string) (endIndex int, command CommandDetails) {
	index := strings.Index(input, CLRF)
	if (index == -1) {
		endIndex = -1
		return endIndex, command
	}

	commLen, _ := strconv.Atoi(input[1:index])
	endIndex = findOccurance(input, "\r",commLen *2)
	if(endIndex == -1) {
		endIndex = -1
		return endIndex, command
	}

	stringToDo := input[:endIndex]
	parts := strings.Split(stringToDo, CLRF)
	temp := []string{}
		for i:=2;i<len(parts);i+=2 {
			temp = append(temp, parts[i])
		}
	

	
	command = CommandDetails{
		command: temp,
		length: len(input[:endIndex+2]),
		raw: input[:endIndex+2],
	}

	return endIndex+2, command
}



func findOccurance(input string, delimiter string, occNum int) int {
		re := regexp.MustCompile(delimiter)
	
		indexes := re.FindAllStringIndex(input, -1)

		if(len(indexes) < occNum) {
			return -1
		}
	
		return indexes[occNum][0]
	
}
