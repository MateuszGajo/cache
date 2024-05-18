package main

import (
	"fmt"
	"reflect"
	"testing"
)

// func TestSingleCommand(t *testing.T) {
// 	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF

// 	res, err := splitMultipleCommandString(tmpString)

// 	if err != nil {
// 		t.Fatalf("got error, %v", err)
// 	}

// 	expected := [][]string{
// 		{"set", "abc", "def"},
// 	}

// 	if !reflect.DeepEqual(res, expected) {
// 		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
// 	}
// }

// func TestSingleCommandSimpleBulk(t *testing.T) {
// 	tmpString := "$3"+CLRF+"abcd"+CLRF

// 	res, err := splitMultipleCommandString(tmpString)

// 	if err != nil {
// 		t.Fatalf("got error, %v", err)
// 	}

// 	expected := [][]string{
// 		{"abcd"},
// 	}

// 	if !reflect.DeepEqual(res, expected) {
// 		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
// 	}
// }

// func TestSingleCommandSimple(t *testing.T) {
// 	tmpString := "+OK"+CLRF

// 	res, err := splitMultipleCommandString(tmpString)

// 	if err != nil {
// 		t.Fatalf("got error, %v", err)
// 	}

// 	expected := [][]string{
// 		{"OK"},
// 	}

// 	if !reflect.DeepEqual(res, expected) {
// 		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
// 	}
// }

// func TestMultipleCommand(t *testing.T) {
// 	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF + "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF

// 	res, err := splitMultipleCommandString(tmpString)

// 	if err != nil {
// 		t.Fatalf("got error, %v", err)
// 	}

// 	expected := [][]string{
// 		{"set", "abc", "def"},
// 		{"set",  "abc", "def"},
// 	}

// 	if !reflect.DeepEqual(res, expected) {
// 		t.Fatalf("wrong answer we expected:%v, and got:%v", expected, res)
// 	}
// }

// func TestInvalid(t *testing.T) {
// 	tmpString := "*3"+CLRF+"$3"+CLRF+"set"+CLRF+"$3"+CLRF+"abc"+CLRF+"$3"+CLRF+"def"+CLRF+"fdfd"+CLRF

// 	_, err := splitMultipleCommandString(tmpString)

// 	if err == nil {
// 		t.Fatal("There should be an error")
// 	}
// }

// func TestEmpty(t *testing.T) {
// 	tmpString := ""

// 	_, err := splitMultipleCommandString(tmpString)

// 	fmt.Print(err)

// 	if err != nil {
// 		t.Fatalf("There is error %v", err)
// 	}
// }

// func TestEmptyaa(t *testing.T) {
// 	tmpString := fmt.Sprintf("*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF,CLRF,CLRF,CLRF,CLRF,CLRF, CLRF)

// 	resp, _ := splitMultipleCommandString2(tmpString)

// 	fmt.Print(resp)

// 	if !reflect.DeepEqual(resp[0], []string{"set", "abc", "def"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"set", "abc", "def"}, resp[0])
// 	}

// }

// func TestEmptyab(t *testing.T) {
// 	tmpString := fmt.Sprintf("+ok%v*3%v$3%vset%v$3%vabc%v$3%vdef%v", CLRF,CLRF,CLRF,CLRF,CLRF,CLRF,CLRF, CLRF)

// 	resp, _ := splitMultipleCommandString2(tmpString)

// 	fmt.Print(resp)

// 	if !reflect.DeepEqual(resp[0], []string{"ok"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"set", "abc", "def"}, resp[0])
// 	}

// 	if !reflect.DeepEqual(resp[1], []string{"set", "abc", "def"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"set", "abc", "def"}, resp[1])
// 	}

// }

// func TestEmptyabc(t *testing.T) {
// 	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2", CLRF, CLRF)

// 	resp, _ := splitMultipleCommandString2(tmpString)

// 	fmt.Print(resp)

// 	if !reflect.DeepEqual(resp[0], []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"}, resp[0])
// 	}

// 	if !reflect.DeepEqual(resp[1], []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
// 		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1])
// 	}

// }

// func TestEmptyabc(t *testing.T) {
// 	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\n123\r\n", CLRF, CLRF)

// 	resp, _ := splitMultipleCommandString2(tmpString)

// 	fmt.Print(resp)

// 	if !reflect.DeepEqual(resp[0], []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"} ){
// 		t.Fatalf("expected %v, got: %v",  []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"}, resp[0])
// 	}

// 	if !reflect.DeepEqual(resp[1], []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
// 		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1])
// 	}

// 	if !reflect.DeepEqual(resp[2], []string{"set", "foo", "123"} ){
// 		t.Fatalf("expected %q, got: %q",  []string{"set", "foo", "123"} , resp[2])
// 	}

// }

func TestEmptyabc(t *testing.T) {
	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$3\r\nset\r\n$3\r\nfoo", CLRF, CLRF)
	

	resp, rest, _ := splitMultipleCommandString2(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0], []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"} ){
		t.Fatalf("expected %v, got: %v",  []string{"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0"}, resp[0])
	}

	if !reflect.DeepEqual(resp[1], []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"} ){
		t.Fatalf("expected %q, got: %q",  []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1])
	}

	if(len(resp) >2) {
		t.Fatalf("should return only two command, returned: %q", resp)
	}

	if !(rest == "*3\r\n$3\r\nset\r\n$3\r\nfoo"){
		t.Fatalf("expected %q, got: %q",  "*3\r\n$3\r\nset\r\n$3\r\nfoo" , rest)
	}

}