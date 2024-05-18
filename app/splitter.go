package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)


func readSyncResp(conn net.Conn) {
	  reader := bufio.NewReader(conn)

	  read, err:= reader.ReadBytes('\n')
	  if err != nil {
		fmt.Printf("Problem reading fullresp response, err:%v", err)
	  }
	  fmt.Print("read untill $", string(read))

	  input2, err:= reader.ReadBytes('\n')
	  if err != nil {
		fmt.Printf("Problem reading szie of rdb file, err:%v", err)
	  }

	  fmt.Println("read until n")
	  fmt.Println(string(input2))
	  fmt.Println(input2)
	  fmt.Println("end  n read")


	  RdbSize, err := strconv.Atoi(strings.Replace(string(input2[1:]), CLRF, "", -1))
	  
	  fmt.Println("value converted")
	  fmt.Println(RdbSize)
	  fmt.Println("end of value")
	  if err != nil {
		fmt.Printf("Read not valid length for rdb file, err: %v", RdbSize)
	  }
	  buff := make([]byte, RdbSize)

	  nof, err := reader.Read(buff)


	  fmt.Println("nof bytes")
	  fmt.Println(nof)
	  fmt.Println("what did we read further")
	  fmt.Println(string(buff[:RdbSize]))
	  fmt.Println(buff[:RdbSize])
	  fmt.Println("what did we read further end")

	  if err != nil {
		fmt.Printf("Problem reading rdb file, err: %v", RdbSize)
	  }
	  fmt.Println("what did you read")
	  fmt.Println(buff[:RdbSize])
	  fmt.Println("end")
	
}

func splitMultipleCommandString2(input string) (res [][]string, tempInp string, err error) {
	tempInp = input
	commands := []string {}
	
	delimiterLen := len(CLRF)

	loop: for (len(tempInp) > 0){
		indicator := tempInp[0]
		fmt.Println("wchodzimy!!!")
		fmt.Printf("%q",tempInp)

	switch (indicator) {
	case 43:
			fmt.Println("Wchodzimy 1")
			index := strings.Index(tempInp, CLRF)
			if(index == -1) {
				break;
			}
			endIndex:= index+delimiterLen
			commands = []string{ tempInp[1:index]}
			res= append(res, commands)
			tempInp = tempInp[endIndex:]
	case 36:
		fmt.Println("Wchodzimy 2")
		index := strings.Index(tempInp, CLRF)
		if(index == -1) {
			break loop;
		}
		commLen, _ := strconv.Atoi(tempInp[1:index])
		endIndex := commLen + index + delimiterLen
		if(len(tempInp) >= endIndex + delimiterLen && tempInp[endIndex +1] == 13 && tempInp[endIndex +2] == 10) {
			//normal command
			commands = []string{ tempInp[index+delimiterLen:endIndex]}
			res= append(res, commands)
			endIndex +=2
			tempInp = tempInp[endIndex:]
		} else {
			//rdb file
			commands = []string{ "rdb-file", tempInp[index+ delimiterLen:endIndex]}
			fmt.Println(commLen, endIndex)
			fmt.Println(len(tempInp))
			res= append(res, commands)
			tempInp = tempInp[endIndex:]
		}

	case 42:
		fmt.Println("Wchodziy?")
		fmt.Println(res)
		index := strings.Index(tempInp, CLRF)
		fmt.Println("inex", index)
		if (index == -1) {
			break loop;
		}
		commLen, _ := strconv.Atoi(tempInp[1:index])
		endIndex := founcOCcurangce(tempInp, commLen *2)
		fmt.Println("next length", commLen)
		fmt.Println("next end index", endIndex)
		if(endIndex == -1) {
			break loop;
		}
		stringToDo := tempInp[:endIndex]
		parts := strings.Split(stringToDo, CLRF)
		temp := []string{}
			for i:=2;i<len(parts);i+=2 {
				temp = append(temp, parts[i])
			}
		res= append(res, temp)
		tempInp = tempInp[endIndex+2:]
	default:
		fmt.Println("unsupported break")
		break loop;
	}
}


	return res,tempInp, nil
}

func founcOCcurangce(input string, occNum int) int {
		fmt.Println("wchodzimy do occurance")
		fmt.Printf("%q\n",input)
		fmt.Println(occNum)
	
		// Regular expression to find all occurrences of "$"
		re := regexp.MustCompile(`\r`)
	
		// Find all indexes of matches
		indexes := re.FindAllStringIndex(input, -1)

		fmt.Println(indexes)

		if(len(indexes) < occNum) {
			return -1
		}
	

		return indexes[occNum][0]
	
}

func splitMultipleCommandString(input string) (res [][]string, err error) {
	var index int
	parts := strings.Split(input, CLRF)
	partsLen := len(parts)
	for len(parts) -1 > index {

		switch parts[index][0] {
		case 43:
			endIndex := index+1
			if(endIndex >= partsLen) {
				fmt.Print("invalid herre")
				return nil, errors.New("look like invalid input")
			}
			command:=parts[index: endIndex][0]
			command = strings.Trim(command, "+")
			temp := []string { command}
			res = append(res, temp)
			index = endIndex
		case 36:
			endIndex := index+2

			if(endIndex > partsLen) {
				return nil, errors.New("look like invalid input")
			}
			res = append(res, parts[index+1: endIndex])
			index = endIndex
		case 42:
			val, _:= strconv.ParseInt(string(parts[index][1]), 10, 32)
			endIndex := index+int(val)*2+1
			if(endIndex >= partsLen) {
				return nil, errors.New("look like invalid input")
			}
			temp := []string{}
			for i:=index+2;i<endIndex;i+=2 {
				temp = append(temp, parts[i])
			}
			res = append(res, temp)
			index = endIndex
		default:
			return nil, errors.New("look like unsporrted resp: "+ string(parts[index][0]))
		}
		
	}

	return res, nil
}
