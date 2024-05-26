package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"testing"
	"time"
)

var _ = func() bool {
    testing.Init()
    return true
}()

type cmdArg []struct {
	args []string
}


//Write this test

func TestG(t *testing.T) {

	fmt.Print("start")
	//This
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

	// 	fmt.Print("hello?")
	// conn, err := net.Dial("tcp", "127.0.0.1:6379")
	


	// conn.Write([]byte("*3\r\n$3\r\nset\r\n$3\r\nabc\r\n$3\r\ndef\r\n"))
	// conn.Write([]byte("*3\r\n$4\r\nwait\r\n$2\r\n1\r\n$3\r\n300\r\n"))

	// buffer := make([]byte, 1000)

	// n, err := conn.Read(buffer)

	// if err != nil {
	// 	fmt.Print("error")
	// 	return;
	// }

	// t.Errorf("value, %v",string(buffer[:n]))

	// // t.Cleanup(func() {
	// // 	// conn.Close()
	// // 	fmt.Print("lets kill process")
	// // 	if err := command.Process.Kill(); err != nil {
	// // 		log.Fatal("failed to kill process: ", err)
	// // 	}
	// // })
  
    // if err != nil {
	// 	// fmt.Print("lets kill process")
	// 	// if err := command.Process.Kill(); err != nil {
	// 	// 	log.Fatal("failed to kill process: ", err)
	// 	// }
    //     t.Error("could not connect to server: ", err)
	// 	return
    // }




}


func TestG2(t *testing.T) {

	fmt.Print("start")
	//This
	commands :=  cmdArg{
        {args: []string{"run",".","--port", "6379"}},
        // {args: []string{"run",".","--port", "6380", "--replicaof","127.0.0.1","6379"}},
    }
	var command *exec.Cmd
	go func(args []string) {
		command = exec.Command("go", args...)
		stdout, err := command.StdoutPipe()

		if err !=nil {
			fmt.Printf("Comgfdrvmand failed: %v", err)
		}

		if err := command.Start(); err != nil {
			fmt.Printf("Comgfdrvmand failed: %v", err)
		}
		in := bufio.NewScanner(stdout)

		for in.Scan() {
			log.Printf(in.Text()) // write each line to your log, or anything you need
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}

		
	}(commands[0].args)

	// go func(args []string) {
	// 	command := exec.Command("go", args...)
	// 	if err := command.Start(); err != nil {
	// 		fmt.Printf("Comgfdrvmand failed: %v", err)
	// 	}
	// }(commands[1].args)


		time.Sleep(1 * time.Second)

	// 	fmt.Print("hello?")
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	


	_, err = conn.Write([]byte("*5\r\n$4\r\nxadd\r\n$3\r\nkey\r\n$3\r\n1-1\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
		if(err != nil ){
			t.Errorf("got err: %q", err)
		}
	// XADD some_key 1-1 foo bar
	// conn.Write([]byte("*3\r\n$4\r\nwait\r\n$2\r\n1\r\n$3\r\n300\r\n"))

	buffer := make([]byte, 1000)

	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Print("error")
		return;
	}

	resp := string(buffer[:n])

	t.Errorf("expected: %q, got:%q", "$3\r\n1-1\r\n", resp)

	if(resp != "$3\r\n1-1\r\n") {
		t.Errorf("expected: %q, got:%q", "$3\r\n1-1\r\n", resp)
	}

	t.Cleanup(func() {
		conn.Close()
		fmt.Print("lets kill process")
		if err := command.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}

	})
  
}