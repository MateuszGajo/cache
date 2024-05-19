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
			if strings.Contains(tempInp, "FULLRESYNC") {
				commands := CommandDetails{
					command: strings.Split(tempInp[1:index], " "),
					length: len(tempInp[:endIndex]),
					raw: tempInp[endIndex:],
				}
				res= append(res, commands)
			}else {
				commands := CommandDetails{
					command: []string{tempInp[1:index]},
					length: len(tempInp[:endIndex]),
					raw: tempInp[:endIndex],
				}
				res= append(res, commands)
			}
			
		
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
			endIndex +=2
			
			commands := CommandDetails{
				command: []string{tempInp[index+delimiterLen:endIndex]},
				length: len(tempInp[:endIndex]),
				raw: tempInp[:endIndex],
			}
			res= append(res, commands)
			tempInp = tempInp[endIndex:]
		} else {
			//rdb file
			fmt.Println(commLen, endIndex)
			fmt.Println(len(tempInp))
			

			commands := CommandDetails{
				command: []string{"rdb-file", tempInp[index+ delimiterLen:endIndex]},
				length: len(tempInp[:endIndex]),
				raw: tempInp[:endIndex],
			}
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
		endIndex := findOccurance(tempInp, "\r",commLen *2)
		if(endIndex == -1) {
			break loop;
		}
		stringToDo := tempInp[:endIndex]
		parts := strings.Split(stringToDo, CLRF)
		temp := []string{}
			for i:=2;i<len(parts);i+=2 {
				temp = append(temp, parts[i])
			}
		

		
		commands := CommandDetails{
			command: temp,
			length: len(tempInp[:endIndex+2]),
			raw: tempInp[:endIndex+2],
		}
		res= append(res, commands)
		tempInp = tempInp[endIndex+2:]
	default:
		fmt.Println("unsupported break")
		break loop;
	}
}


	return res,tempInp, nil
}

func findOccurance(input string, delimiter string, occNum int) int {
		re := regexp.MustCompile(delimiter)
	
		indexes := re.FindAllStringIndex(input, -1)

		if(len(indexes) < occNum) {
			return -1
		}
	
		return indexes[occNum][0]
	
}
