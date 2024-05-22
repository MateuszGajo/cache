package main

import (
	"fmt"
	"net"
	"testing"
)

var _ = func() bool {
    testing.Init()
    return true
}()

type cmdArg []struct {
	args []string
}


func TestG(t *testing.T) {

	fmt.Print("start")
	// commands :=  cmdArg{
    //     {args: []string{"run",".","--port", "6379"}},
    //     // {args: []string{"run",".","--port", "6380", "--replicaof","127.0.0.1","6379"}},
    // }
	// var command *exec.Cmd

	// go func(args []string) {
	// 	command = exec.Command("go", args...)
	// 	fmt.Print("start")
	// 	if err := command.Start(); err != nil {
	// 		fmt.Printf("Comgfdrvmand failed: %v", err)
	// 	}
		
	// }(commands[0].args)

	// go func(args []string) {
	// 	command := exec.Command("go", args...)
	// 	if err := command.Start(); err != nil {
	// 		fmt.Printf("Comgfdrvmand failed: %v", err)
	// 	}
	// }(commands[1].args)


		// time.Sleep(1 * time.Second)

		fmt.Print("hello?")
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	
    // defer func(conn net.Conn){
	// 	// conn.Close()
	// 	fmt.Print("lets kill process")
	// 	if err := command.Process.Kill(); err != nil {
	// 		log.Fatal("failed to kill process: ", err)
	// 	}
	// }(conn)

	conn.Write([]byte("*3\r\n$3\r\nset\r\n$3\r\nabc\r\n$3\r\ndef\r\n"))
	conn.Write([]byte("*3\r\n$4\r\nwait\r\n$1\r\n1\r\n$3\r\n300\r\n"))

	buffer := make([]byte, 1000)

	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Print("error")
		return;
	}

	t.Errorf("value, %v",string(buffer[:n]))

	// t.Cleanup(func() {
	// 	// conn.Close()
	// 	fmt.Print("lets kill process")
	// 	if err := command.Process.Kill(); err != nil {
	// 		log.Fatal("failed to kill process: ", err)
	// 	}
	// })
  
    if err != nil {
		// fmt.Print("lets kill process")
		// if err := command.Process.Kill(); err != nil {
		// 	log.Fatal("failed to kill process: ", err)
		// }
        t.Error("could not connect to server: ", err)
		return
    }




}