package main

import (
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

func handleSet(key, value string, expiryTime *int) string {
	lock.Lock()
	defer lock.Unlock()
	if expiryTime != nil {
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Now().Add(time.Duration(*expiryTime) * time.Millisecond),
		}
	} else {
		m[key] = CustomSetStore{
			Value:    value,
			ExpireAt: time.Time{},
		}
	}

	return "+OK\r\n"
}


func handleGet(key string) string {
	defer lock.RUnlock()
	lock.RLock()
	r, ok := m[key]
	if !ok || (time.Now().After(r.ExpireAt) && r.ExpireAt != time.Time{}) {
		return "$-1\r\n"
	}

	return BuildResponse(r.Value)
}