package main

import (
	"fmt"
	"strings"
)

func RESPSimpleString(input string) string {
	arg := strings.Split(input, CLRF)[0]
	arg = strings.Trim(arg, "+")

	return arg
}

func RESPArray(input string) []string {
	args := strings.Split(input, CLRF)

	fmt.Println(args)

	command := make([]string, 0, ((len(args) - 1) / 2))

	for i := 2; i < len(args); i = i + 2 {
		command = append(command, args[i])
	}

	return command
}

func RESPBulkString(input string) []string {
	args := strings.Split(input, CLRF)

	fmt.Println(args)

	command := make([]string, 0, ((len(args) - 1) / 2))

	for i := 2; i < len(args); i = i + 2 {
		command = append(command, args[i])
	}

	return command
}

func BuildSimpleString(message string) string {
	fmt.Print("Hello")
	return fmt.Sprintf("+%v%v", message, CLRF)
}

func BuildBulkString(message string) string {
	return fmt.Sprintf("$%v%v%s%v", len(message), CLRF, message, CLRF)
}

func BuildNullBulkString() string {
	return fmt.Sprintf("$-1%v", CLRF)
}

func BuildRespInt(val int) string {
	return fmt.Sprintf(":%v%v", val, CLRF)
}

func BuildSimpleError(errType string, message string) string {
	return fmt.Sprintf("-%v %v", errType, message)
}

func BuildRESPArray(args []string) string {
	resp := fmt.Sprintf("*%v%v", len(args), CLRF)

	for _, arg := range args {
		resp += fmt.Sprintf("$%v%v%v%v", len(arg), CLRF, arg, CLRF)
	}

	return resp
}
   