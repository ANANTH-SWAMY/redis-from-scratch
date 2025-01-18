package main

import(
	"sync"
)

type storeValue struct {
	Bulk string
	HashStore map[string]string
	IsHash bool
}

var store = make(map[string]storeValue)
var storeMu = sync.RWMutex{}
