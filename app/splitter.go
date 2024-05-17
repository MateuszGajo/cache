package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)


func readSyncResp(conn net.Conn) {
	  reader := bufio.NewReader(conn)

	  read, err:= reader.ReadBytes('$')
	  if err != nil {
		fmt.Printf("Problem reading fullresp response, err:%v", err)
	  }
	  fmt.Print("read untill $", string(read))

	  input2, err:= reader.ReadBytes('n')
	  if err != nil {
		fmt.Printf("Problem reading szie of rdb file, err:%v", err)
	  }

	  fmt.Print("read until n", input2)


	  RdbSize, err := strconv.Atoi(strings.Replace(string(input2), CLRF, "", -1))
	  if err != nil {
		fmt.Printf("Read not valid length for rdb file, err: %v", RdbSize)
	  }
	  buff := make([]byte, RdbSize)

	  _, err = reader.Read(buff)

	  if err != nil {
		fmt.Printf("Problem reading rdb file, err: %v", RdbSize)
	  }
	  fmt.Println("what did you read")
	  fmt.Println(buff[:RdbSize])
	  fmt.Println("end")
	
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
