package main

var handlers = map[string] func([]Value) Value {
	"PING": ping,
	"SET": set,
	"GET": get,
}

var store = make(map[string]string)

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

func set(args []Value) Value {
	if len(args) != 2 {
		v := Value{
			typ: "error",
			str: "ERR wrong number of arguments for command",
		}

		return v
	}

	key := args[0].bulk
	value := args[1].bulk

	store[key] = value

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
			str: "ERR wrong number of arguments for command",
		}

		return v
	}

	key := args[0].bulk
	value := store[key]

	if value == "" {
		v := Value{
			typ: "null",
		}

		return v
	}

	v := Value{
		typ: "string",
		str: value,
	}

	return v
}
