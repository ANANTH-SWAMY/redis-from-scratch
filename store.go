package main

import(
	"sync"
)

type storeValue struct {
	bulk string
	hashStore map[string]string
	isHash bool
}

var store = make(map[string]storeValue)
var storeMu = sync.RWMutex{}
