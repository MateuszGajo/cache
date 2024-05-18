package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Aa struct {
	command []string
	length int
	raw string
}


func splitMultipleCommandString(input string) (res []Aa, tempInp string, err error) {
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
				commands := Aa{
					command: strings.Split(tempInp[1:index], " "),
					length: len(tempInp[:endIndex]),
					raw: tempInp[endIndex:],
				}
				res= append(res, commands)
			}else {
				commands := Aa{
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
			
			commands := Aa{
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
			

			commands := Aa{
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
		

		
		commands := Aa{
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
