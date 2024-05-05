package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var m = map[string]CustomSetStore{}
var lock = sync.RWMutex{}

func handleRemove(key string) {
	lock.Lock()
	defer lock.Unlock()
	delete(m, key)
	//TO DO
}

func handleSet(key, value string, expiryTime *int) bool {
	lock.Lock()
	defer lock.Unlock()
	if expiryTime != nil {
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Now().Add(time.Duration(*expiryTime) * time.Millisecond),
		}
	} else {
		fmt.Print("save")
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Time{},
		}
	}

	return true
}


func handleGet(key string) (string, error) {
	defer lock.RUnlock()
	lock.RLock()
	r, ok := m[key]
	if !ok {
		return "", errors.New("Problem while getting value")
	}
	if (time.Now().After(r.ExpireAt) && r.ExpireAt != time.Time{}) {
		return "", nil
	}

	return r.Value, nil
}