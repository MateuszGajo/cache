package main

import (
	"testing"
	"time"
)

func TestReadSetValue(t *testing.T) {
	key := "abc"
	value := "123"
	res := handleSet(key, value, nil, "string")

	if !res {
		t.Fatalf("problem setting %v with %v value", key, value)
	}

	readVal, _ := HandleGetString(key, DatabaseConfig{})

	if readVal != "123" {
		t.Fatalf("values are diffrent, set: %v, get:%v", value, readVal)
	}
}

func TestReadNotExistValue(t *testing.T) {
	key := "notExist"
	readVal, err := HandleGetString(key, DatabaseConfig{})

	if err == nil {
		t.Fatalf("Value should be empty, isnted recived: %v, err:%v", readVal, err)
	}
}

func TestExpiredValueShouldReturnEmpty(t *testing.T) {
	key := "abc"
	value := "123"
	expiry := 100
	res := handleSet(key, value, &expiry, "string")

	if !res {
		t.Fatalf("problem setting %v with %v value", key, value)
	}

	time.Sleep(time.Millisecond * time.Duration(expiry))

	readVal, err := HandleGetString(key, DatabaseConfig{})

	if err == nil {
		t.Fatalf("Value should be empty, isnted recived: %v, err:%v", readVal, err)
	}
}

func TestExpiredValueShouldReturnEmpty1(t *testing.T) {
	handleSet("abc", "123", nil, "string")
	handleSet("def", Stream{}, nil, "stream")

	entryType := HandleGetType("abc")

	if entryType != "string" {
		t.Fatalf("for key abc and value 123 type should be string, but recived: %v", entryType)
	}

	entryType = HandleGetType("def")

	if entryType != "stream" {
		t.Fatalf("for key abc and value 123 type should be string, but recived: %v", entryType)
	}

}
