package main

import (
	"fmt"
)

var handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"DEL":     del,
	"MGET":    mget,
	"MSET":    mset,
	"HSET":    hset,
	"HGET":    hget,
	"HDEL":    hdel,
	"EXISTS":  exists,
	"HEXISTS": hexists,
	"COMMAND": command,
}

func wrongNoOfArguments(cmd string) Value {
	v := Value{
		typ: "error",
		str: fmt.Sprintf("ERR wrong number of arguments for '%v' command", cmd),
	}

	return v
}

func unknownCommand(cmd string) Value {
	v := Value{
		typ: "error",
		str: fmt.Sprintf("ERR unknown command '%v'", cmd),
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

func set(args []Value) Value {
	if len(args) != 2 {
		return wrongNoOfArguments("set")
	}

	key := args[0].bulk
	value := args[1].bulk

	storeMu.Lock()
	store[key] = storeValue{Bulk: value, IsHash: false}
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
		store[key] = storeValue{Bulk: value, IsHash: false}
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
	value, ok := store[key]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "null",
		}

		return v
	}

	if value.IsHash{
		v := Value{
			typ: "error",
			str: "WRONGTYPE Operation against a key holding the wrong kind of value",
		}

		return v
	}	

	v := Value{
		typ: "bulk",
		bulk: value.Bulk,
	}

	return v
}

func mget(args []Value) Value {
	if len(args) == 0 {
		return wrongNoOfArguments("mget")
	}

	v := Value{
		typ: "array",
		array: make([]Value, 0),
	}

	for i := 0; i < len(args); i++ {
		key := args[i].bulk

		storeMu.RLock()
		value, ok := store[key]
		storeMu.RUnlock()

		if ok && !value.IsHash {
			newElement := Value{
				typ: "bulk",
				bulk: value.Bulk,
			}

			v.array = append(v.array, newElement)
		} else {
			newElement := Value{
				typ: "null",
			}

			v.array = append(v.array, newElement)
		}
	}

	return v
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

func hset(args []Value) Value {
	if len(args) < 3 {
		return wrongNoOfArguments("hset")
	}

	key := args[0].bulk
	args = args[1:]

	if len(args) % 2 != 0 {
		return wrongNoOfArguments("hset")
	}

	storeMu.RLock()
	_, ok := store[key]
	storeMu.RUnlock()

	if !ok {
		storeMu.Lock()

		store[key] = storeValue{
			HashStore: make(map[string]string),
			IsHash: true,
		}

		storeMu.Unlock()
	}

	if store[key].IsHash == false {
		v := Value{
			typ: "error",
			str: "WRONGTYPE Operation against a key holding the wrong kind of value",
		}

		return v
	}

	count := 0

	for i := 0; i < len(args); i = i + 2 {
		field := args[i].bulk
		value := args[i+1].bulk

		storeMu.RLock()
		_, ok := store[key].HashStore[field]
		storeMu.RUnlock()

		storeMu.Lock()
		store[key].HashStore[field] = value
		storeMu.Unlock()

		if !ok {
			count ++
		}
	}

	v := Value{
		typ: "integer",
		integer: count,
	}

	return v
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return wrongNoOfArguments("hget")
	}

	key := args[0].bulk
	field := args[1].bulk

	storeMu.RLock()
	HashStore, ok := store[key]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "null",
		}

		return v
	}

	if !HashStore.IsHash {
		v := Value{
			typ: "error",
			str: "WRONGTYPE Operation against a key holding the wrong kind of value",
		}

		return v
	}

	storeMu.RLock()
	val, ok := HashStore.HashStore[field]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "null",
		}

		return v
	}

	v := Value{
		typ: "bulk",
		bulk: val,
	}

	return v
}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return wrongNoOfArguments("hdel")
	}

	key := args[0].bulk
	fields := args[1:]

	storeMu.RLock()
	HashStore, ok := store[key]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "integer",
			integer: 0,
		}

		return v
	}

	if !HashStore.IsHash {
		v := Value{
			typ: "error",
			str: "WRONGTYPE Operation against a key holding the wrong kind of value",
		}

		return v
	}

	count := 0

	for i := 0; i < len(fields); i++ {
		field := fields[i].bulk

		storeMu.RLock()
		_, ok := store[key].HashStore[field]
		storeMu.RUnlock()

		if ok {
			storeMu.Lock()
			delete(store[key].HashStore, field)
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

func hexists(args []Value) Value {
	if len(args) != 2 {
		return wrongNoOfArguments("hexists")
	}

	key := args[0].bulk
	field := args[1].bulk

	storeMu.RLock()
	HashStore, ok := store[key]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "integer",
			integer: 0,
		}

		return v
	}

	if !HashStore.IsHash {
		v := Value{
			typ: "error",
			str: "WRONGTYPE Operation against a key holding the wrong kind of value",
		}

		return v
	}

	storeMu.RLock()
	_, ok = HashStore.HashStore[field]
	storeMu.RUnlock()

	if !ok {
		v := Value{
			typ: "integer",
			integer: 0,
		}

		return v
	}

	v := Value{
		typ: "integer",
		integer: 1,
	}

	return v
}

// Placeholder to ignore the initial 'COMMAND DOCS' command sent by redis-cli
func command(args []Value) Value {
	v := Value{
		typ: "string",
		str: "OK",
	}

	return v
}
