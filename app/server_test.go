package main

import (
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

	// fmt.Print("start")
	// commands :=  cmdArg{
    //     {args: []string{"run",".","--port", "6379"}},
    //     {args: []string{"run",".","--port", "6380", "--replicaof","127.0.0.1","6379"}},
    // }

	// go func(args []string) {
	// 	command := exec.Command("go", args...)
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



  

	// 	fmt.Print("hello?")
	// conn, err := net.Dial("tcp", "127.0.0.1:6379")
    // if err != nil {
    //     t.Error("could not connect to server: ", err)
    // }
    // defer conn.Close()

	// fmt.Print("hello?")
	// conn1, err := net.Dial("tcp", "127.0.0.1:6380")
    // if err != nil {
    //     t.Error("could not connect to server: ", err)
    // }
    // defer conn.Close()



	// conn.Write([]byte(fmt.Sprintf("*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF, CLRF, CLRF, CLRF, CLRF, CLRF, CLRF)))
	// buf := make([]byte, 1024)
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	 t.Error("could not connect to server: ", err)
	// }

	// input :=string(buf[:n])

	// fmt.Print(input)
	// if input != "+OK" + CLRF {
	// 	t.Errorf("Wront output, expected %v, got %v ","+Ok" + CLRF, input)
	// }

	// conn.Write([]byte(fmt.Sprintf("*3%v$3%vset%v$3%vdef%v$3%v123%v", CLRF, CLRF, CLRF, CLRF, CLRF, CLRF, CLRF)))
	// buf = make([]byte, 1024)
	// n, err = conn.Read(buf)
	// if err != nil {
	// 	 t.Error("could not connect to server: ", err)
	// }

	// input =string(buf[:n])

	// fmt.Print(input)
	// if input != "+OK" + CLRF {
	// 	t.Errorf("Wront output, expected %v, got %v ","+Ok" + CLRF, input)
	// }

	// conn.Write([]byte(fmt.Sprintf("*3%v$3%vset%v$3%vbar%v$3%vbvbvv%v", CLRF, CLRF, CLRF, CLRF, CLRF, CLRF, CLRF)))
	// buf = make([]byte, 1024)
	// n, err = conn.Read(buf)
	// if err != nil {
	// 	 t.Error("could not connect to server: ", err)
	// }

	// input =string(buf[:n])

	// fmt.Print(input)
	// if input != "+OK" + CLRF {
	// 	t.Errorf("Wront output, expected %v, got %v ","+Ok" + CLRF, input)
	// }

	// conn1.Write([]byte(fmt.Sprintf("*2%v$3%vget%v$3%vbar%v", CLRF, CLRF, CLRF, CLRF, CLRF)))
	// buf = make([]byte, 1024)
	// n, err = conn1.Read(buf)
	// if err != nil {
	// 	 t.Error("could not connect to server: ", err)
	// }

	// input =string(buf[:n])

	// fmt.Print(input)
	// if input != "$5"+CLRF+"bvbvv"+CLRF {
	// 	t.Errorf("Wront output, expected %v, got %v ","$5"+CLRF+"bvbvv"+CLRF, input)
	// }

	// conn1.Write([]byte(fmt.Sprintf("*2%v$3%vget%v$3%vdef%v", CLRF, CLRF, CLRF, CLRF, CLRF)))
	// buf = make([]byte, 1024)
	// n, err = conn1.Read(buf)
	// if err != nil {
	// 	 t.Error("could not connect to server: ", err)
	// }

	// input =string(buf[:n])

	// fmt.Print(input)
	// if input != "$3\\r\\n123\\r\\n" {
	// 	t.Errorf("Wront output, expected %v, got %v ","$3\\r\\n123\\r\\n", input)
	// }

	//test replicas

}