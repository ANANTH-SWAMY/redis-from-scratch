package main

import (
	"sync"
)

var handlers = map[string] func([]Value) Value {
	"PING": ping,
	"SET": set,
	"GET": get,
	"DEL": del,
	"EXISTS": exists,
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
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'ping' command",
		}

		return v
	}

	v := Value{
		typ: "bulk",
		bulk: args[0].bulk,
	}

	return v
}

var store = make(map[string]string)
var storeMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'set' command",
		}

		return v
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

func get(args []Value) Value {
	if len(args) != 1 {
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'get' command",
		}

		return v
	}

	key := args[0].bulk

	storeMu.RLock()
	value, ok := store[key]
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

func exists(args []Value) Value {
	if len(args) == 0 {
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'exists' command",
		}

		return v
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

func del(args []Value) Value {
	if len(args) == 0 {
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'del' command",
		}

		return v
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
