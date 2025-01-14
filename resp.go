package main

import (
	"bufio"
	"strconv"
)

const (
	STRING = '+'
	ERROR = '-'
	INTEGER = ':'
	BULK = '$'
	ARRAY = '*'
	NULL = '_'
)

type Value struct {
	typ      string
	str      string
	integer  int
	bulk     string
	array    []Value
}

func readLine(r *bufio.Reader) ([]byte, int, error) {
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

func readInteger(r *bufio.Reader) (int, error) {
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

func readArray(r *bufio.Reader) (Value, error) {
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

func readBulk(r *bufio.Reader) (Value, error) { 
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

func parse(r *bufio.Reader) (Value, error) {
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

func writeInteger(v Value) []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.integer)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func writeArray(v Value) []byte {
	bytes := make([]byte, 0)

	array := v.array

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(array))...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len(array); i++ {
		element := array[i]
		bytesToBeAppended := make([]byte, 0)

		switch element.typ {
		case "string":
			bytesToBeAppended = writeString(element)

		case "bulk":
			bytesToBeAppended = writeBulk(element)

		case "integer":
			bytesToBeAppended = writeInteger(element)
	
		case "null":
			bytesToBeAppended = writeNull()

		case "array":
			bytesToBeAppended = writeArray(element)
		}

		bytes = append(bytes, bytesToBeAppended...)
	}

	return bytes
}

func writeError(v Value) []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func writeNull() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, NULL)
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

	case "integer":
		return writeInteger(v)

	case "array":
		return writeArray(v)

	case "error":
		return writeError(v)

	case "null":
		return writeNull()

	default:
		return make([]byte, 0)
	}
}
