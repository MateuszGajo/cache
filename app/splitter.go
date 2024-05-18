package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


func splitMultipleCommandString(input string) (res [][]string, tempInp string, err error) {
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
			commands := []string{ }
			endIndex:= index+delimiterLen
			if strings.Contains(tempInp, "FULLRESYNC") {
				commands = strings.Split(tempInp[1:index], " ")
			}else {
				commands = append(commands, tempInp[1:index])
			}
			
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
			commands := []string{ tempInp[index+delimiterLen:endIndex]}
			res= append(res, commands)
			endIndex +=2
			tempInp = tempInp[endIndex:]
		} else {
			//rdb file
			commands := []string{ "rdb-file", tempInp[index+ delimiterLen:endIndex]}
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
		res= append(res, temp)
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
