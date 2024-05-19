package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSingleCommand(t *testing.T) {
	tmpString := fmt.Sprintf("*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF,CLRF,CLRF,CLRF,CLRF,CLRF, CLRF)

	resp, _,_ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0].command, []string{"set", "abc", "def"} ){
		t.Fatalf("expected %v, got: %v",  []string{"set", "abc", "def"}, resp[0].command)
	}

}

func TestMultipleCommand(t *testing.T) {
	tmpString := fmt.Sprintf("+ok%v*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF,CLRF,CLRF,CLRF,CLRF,CLRF,CLRF, CLRF)

	resp, _,_ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0].command, []string{"ok"} ){
		t.Fatalf("expected %v, got: %v",  []string{"ok"}, resp[0].command)
	}

	if !reflect.DeepEqual(resp[1].command, []string{"set", "abc", "def"} ){
		t.Fatalf("expected %v, got: %v",  []string{"set", "abc", "def"}, resp[1].command)
	}

}

func TestCommandWithRdbFile(t *testing.T) {
	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2", CLRF, CLRF)

	resp, _, _ := splitMultipleCommandString(tmpString)

	if !reflect.DeepEqual(resp[0].command, []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"} ){
		t.Fatalf("expected %v, got: %v with len: %v",  []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}, resp[0].command, len(resp[0].command))
	}

	if !reflect.DeepEqual(resp[1].command, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1].command)
	}

}

func TestWithCommandAfterRdbFile(t *testing.T) {
	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\n123\r\n", CLRF, CLRF)

	resp, _,_ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0].command, []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"} ){
		t.Fatalf("expected %v, got: %v",  []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}, resp[0].command)
	}

	if !reflect.DeepEqual(resp[1].command, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1].command)
	}

	if !reflect.DeepEqual(resp[2].command, []string{"set", "foo", "123"} ){
		t.Fatalf("expected %q, got: %q",  []string{"set", "foo", "123"} , resp[2].command)
	}

}
func TestWithCommandAfterRdbFile2(t *testing.T) {
	tmpString := "$88\r\nREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"

	resp, _,_ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)



	if !reflect.DeepEqual(resp[0].command, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[0].command)
	}

	if !reflect.DeepEqual(resp[1].command, []string{"REPLCONF", "GETACK", "*"} ){
		t.Fatalf("expected %q, got: %q",  []string{"REPLCONF", "GETACK", "*"} , resp[1].command)
	}

}





// func TestBreakComamnd(t *testing.T) {
// 	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$3\r\nset\r\n$3\r\nfoo", CLRF, CLRF)
	

// 	resp, rest, _ := splitMultipleCommandString(tmpString)

// 	fmt.Print(resp)

// 	if !reflect.DeepEqual(resp[0].command, []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}, resp[0].command)
// 	}

// 	if !reflect.DeepEqual(resp[1].command, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
// 		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1].command)
// 	}

// 	if(len(resp) >2) {
// 		t.Fatalf("should return only two command, returned: %q", resp)
// 	}

// 	if !(rest == "*3\r\n$3\r\nset\r\n$3\r\nfoo"){
// 		t.Fatalf("expected %q, got: %q",  "*3\r\n$3\r\nset\r\n$3\r\nfoo" , rest)
// 	}

// }