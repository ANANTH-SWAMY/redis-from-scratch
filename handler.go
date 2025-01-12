package main

var handlers = map[string] func([]Value) Value {
	"PING": ping,
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
			str: "wrong number of arguments for 'ping' command",
		}

		return v
	}

	v := Value{
		typ: "bulk",
		bulk: args[0].bulk,
	}

	return v
}
