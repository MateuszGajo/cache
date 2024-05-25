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

	readVal := handleGet(key)


	if (readVal != CustomSetStore{}) {
		t.Fatalf("values are diffrent, set: %v, get:%v", value, readVal)
	}
}

func TestReadNotExistValue(t *testing.T) {
	key := "notExist"
	readVal := handleGet(key)

	if (readVal != CustomSetStore{}) {
		t.Fatalf("Value should be empty, isnted recived: %v", readVal)
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

	time.Sleep(time.Millisecond * time.Duration(expiry) )

	readVal := handleGet(key)


	if (readVal != CustomSetStore{}) {
		t.Fatalf("Value should be empty, isnted recived: %v", readVal)
	}
}