package main

var handlers = map[string] func([]Value) Value {
	"PING": ping,
}

func ping(args []Value) Value {
	v := Value{
		typ: "string",
		str: "PONG",
	}

	return v
}
