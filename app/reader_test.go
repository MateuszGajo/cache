package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSingleCommand(t *testing.T) {
	tmpString := BuildPrimitiveRESPArray([]string{"set", "abc", "def"})

	resp, _, _ := splitMultipleCommandString(tmpString)
	commandOneArray := resp[0].data.([]string)

	if !reflect.DeepEqual(commandOneArray, []string{"set", "abc", "def"}) {
		t.Fatalf("expected %v, got: %v", []string{"set", "abc", "def"}, commandOneArray)
	}

}

func TestMultipleCommand(t *testing.T) {
	tmpString := BuildSimpleString("ok") + BuildPrimitiveRESPArray([]string{"set", "abc", "def"})

	resp, _, _ := splitMultipleCommandString(tmpString)

	commandOneSimpleString := resp[0].data.(string)
	commandTwoArray := resp[1].data.([]string)

	if !reflect.DeepEqual(commandOneSimpleString, "ok") {
		t.Fatalf("expected %v, got: %v", "ok", commandOneSimpleString)
	}

	if !reflect.DeepEqual(commandTwoArray, []string{"set", "abc", "def"}) {
		t.Fatalf("expected %v, got: %v", []string{"set", "abc", "def"}, commandTwoArray)
	}

}

func TestCommandWithRdbFile(t *testing.T) {
	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2", CLRF, CLRF)

	resp, _, _ := splitMultipleCommandString(tmpString)

	if !reflect.DeepEqual(resp[0].data, []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}) {
		t.Fatalf("expected %v, got: %v with len: %v", []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}, resp[0].data, len(resp[0].data.(string)))
	}

	if !reflect.DeepEqual(resp[1].data, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}) {
		t.Fatalf("expected %q, got: %q", []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1].data)
	}

}

func TestWithCommandAfterRdbFile(t *testing.T) {
	tmpString := fmt.Sprintf("+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0%v$88%vREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\n123\r\n", CLRF, CLRF)

	resp, _, _ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0].data, []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}) {
		t.Fatalf("expected %v, got: %v", []string{"FULLRESYNC", "75cd7bc10c49047e0d163660f3b90625b1af31dc", "0"}, resp[0].data)
	}

	if !reflect.DeepEqual(resp[1].data, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}) {
		t.Fatalf("expected %q, got: %q", []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[1].data)
	}

	if !reflect.DeepEqual(resp[2].data, []string{"set", "foo", "123"}) {
		t.Fatalf("expected %q, got: %q", []string{"set", "foo", "123"}, resp[2].data)
	}

}
func TestWithCommandAfterRdbFile2(t *testing.T) {
	tmpString := "$88\r\nREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"

	resp, _, _ := splitMultipleCommandString(tmpString)

	fmt.Print(resp)

	if !reflect.DeepEqual(resp[0].data, []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}) {
		t.Fatalf("expected %q, got: %q", []string{"rdb-file", "REDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\b\xbce\xfa\bused-mem°\xc4\x10\x00\xfa\baof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ\xa2"}, resp[0].data)
	}

	if !reflect.DeepEqual(resp[1].data, []string{"REPLCONF", "GETACK", "*"}) {
		t.Fatalf("expected %q, got: %q", []string{"REPLCONF", "GETACK", "*"}, resp[1].data)
	}

}
