package main

import (
	"fmt"
	"sync"
)

var handlers = map[string] func([]Value) Value {
	"PING": ping,
	"SET": set,
	"GET": get,
	"DEL": del,
	"MSET": mset,
	"EXISTS": exists,
}

func wrongNoOfArguments(cmd string) Value {
	v := Value{
		typ: "error",
		str: fmt.Sprintf("ERR wrong number of arguments for '%v' command", cmd),
	}

	return v
}

func ping(args []Value) Value {
	if len(args) == 0 {
		v := Value{
			typ: "string",
			str: "PONG",
		}

		return v
	}

	if len(args) > 1 {
		return wrongNoOfArguments("ping")
	}

	v := Value{
		typ: "bulk",
		bulk: args[0].bulk,
	}

	return v
}

var store = make(map[string]any)
var storeMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return wrongNoOfArguments("set")
	}

	key := args[0].bulk
	value := args[1].bulk

	storeMu.Lock()
	store[key] = value
	storeMu.Unlock()

	v := Value{
		typ: "string",
		str: "OK",
	}

	return v
}

func mset(args []Value) Value {
	if len(args) < 2 {
		return wrongNoOfArguments("mset")
	}

	for i := 0; i < len(args); i = i + 2 {
		if i + 1 >= len(args) {
			break
		}

		key := args[i].bulk
		value := args[i+1].bulk

		storeMu.Lock()
		store[key] = value
		storeMu.Unlock()
	}

	v := Value{
		typ: "string",
		str: "OK",
	}

	return v
}

func get(args []Value) Value {
	if len(args) != 1 {
		return wrongNoOfArguments("get")
	}

	key := args[0].bulk

	storeMu.RLock()
	value, ok := store[key].(string)
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "null",
		}

		return v
	}

	v := Value{
		typ: "bulk",
		bulk: value,
	}

	return v
}

func mget(args []Value) Value {
	return Value{}
}

func del(args []Value) Value {
	if len(args) == 0 {
		return wrongNoOfArguments("del")
	}

	count := 0
	for i := 0; i < len(args); i++ {
		key := args[i].bulk

		storeMu.RLock()
		_, ok := store[key]
		storeMu.RUnlock()

		if ok {
			storeMu.Lock()
			delete(store, key)
			storeMu.Unlock()

			count++
		}
	}

	v := Value{
		typ: "integer",
		integer: count,
	}

	return v
}

func exists(args []Value) Value {
	if len(args) == 0 {
		return wrongNoOfArguments("exists")
	}

	count := 0
	for i := 0; i < len(args); i++ {
		key := args[i].bulk

		storeMu.RLock()
		_, ok := store[key]
		storeMu.RUnlock()

		if ok {
			count++
		}
	}

	v := Value{
		typ: "integer",
		integer: count,
	}
	
	return v
}
