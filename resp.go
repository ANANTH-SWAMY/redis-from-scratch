package main

import (
	"strings"
	"strconv"
)

const (
	STRING = '+'
	ERROR = '-'
	INTEGER = ':'
	BULK = '$'
	ARRAY = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

func readLine(r *strings.Reader) ([]byte, int, error) {
	numberBuffer := make([]byte, 0)
	n := 0

	for {
		byt, err := r.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		if string(byt) == "\r"{
			r.ReadByte() // read and discard the \n
			break
		}

		numberBuffer = append(numberBuffer, byt)
		n++
	}

	return numberBuffer, n, nil
}

func readInteger(r *strings.Reader) (int, error) {
	integer, _, err := readLine(r)
	if err != nil {
		return 0, err
	}

	lengthInt, err := strconv.Atoi(string(integer))
	if err != nil {
		return 0, err
	}

	return lengthInt, nil
}

func readArray(r *strings.Reader) (Value, error) {
	v := Value{}
	v.typ = "array"

	count, err := readInteger(r)
	if err != nil {
		return v, err
	}

	v.array = make([]Value, 0)

	for i := 0; i < count; i++ {
		parsed, err := parse(r)
		if err != nil {
			return v, err
		}

		v.array = append(v.array, parsed)
	}

	return v, nil
}

func readBulk(r *strings.Reader) (Value, error) { 
	v := Value{}
	v.typ = "bulk"

	length, err := readInteger(r)
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.Read(bulk)

	v.bulk = string(bulk)

	// read and discard remaining \r\n
	r.ReadByte()
	r.ReadByte()

	return v, nil
}

func parse(r *strings.Reader) (Value, error) {
	dataType, err := r.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case ARRAY:
		return readArray(r)
	case BULK:
		return readBulk(r)
	default:
		return Value{}, nil
	}
}

func writeString(v Value) []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func writeBulk(v Value) []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func writeError(v Value) []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func writeRESP(v Value) []byte {
	dataType := v.typ

	switch dataType {
	case "string": 
		return writeString(v)
	case "bulk":
		return writeBulk(v)
	case "error":
		return writeError(v)
	default:
		return make([]byte, 0)
	}
}
