package main

import (
	"errors"
	"sync"
	"time"
)

type CustomSetStore struct {
	Value    interface{}
	ExpireAt time.Time
	Type     string
}

var m = map[string]CustomSetStore{}
var lock = sync.RWMutex{}

func handleSet(key string, value interface{}, expiryTime *int, recordType string) bool {
	lock.Lock()
	defer lock.Unlock()
	if expiryTime != nil {
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Now().Add(time.Duration(*expiryTime) * time.Millisecond),
			Type:     recordType,
		}
	} else {
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Time{},
			Type:     recordType,
		}
	}

	return true
}

func handleGet(key string) (CustomSetStore, error) {
	defer lock.RUnlock()
	lock.RLock()
	r, ok := m[key]

	if !ok || (time.Now().After(r.ExpireAt) && r.ExpireAt != time.Time{}) {
		return CustomSetStore{}, errors.New("There is no value")
	}

	return r, nil
}

func HandleGetString(key string, dbConfig DatabaseConfig) (string, error) {
	res, err := handleGet(key)

	if (res == CustomSetStore{}) {

		resp := readFile(dbConfig.dirName + "/" + dbConfig.fileName)
		res := ""
		for _, v := range resp {
			if v.key == key && (v.exp.UnixMilli() == 0 || v.exp.After(time.Now())) {
				res = v.value
			}
		}
		if res != "" {
			return res, nil
		}

		return "", err

	}

	stringValue, ok := res.Value.(string)

	if !ok {
		panic("Value isn't type of string")
	}

	return stringValue, err
}

func HandleGetStream(key string) Stream {
	res, err := handleGet(key)

	if (err != nil || res == CustomSetStore{}) {
		return Stream{}
	}

	stringValue, ok := res.Value.(Stream)

	if !ok {
		panic("Value isn't type of Stream")
	}

	return stringValue
}

func HandleGetType(key string) string {
	res, err := handleGet(key)

	if (err != nil || res == CustomSetStore{}) {
		return "none"
	}

	return res.Type
}
