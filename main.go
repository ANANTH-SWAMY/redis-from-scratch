package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on port 6379...")

	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	for {
		buffer := make([]byte, 1024)

		_, err := connection.Read(buffer)
		if err != nil {

			if err != io.EOF {
				fmt.Println("Unable to read:", err)
			}

			fmt.Println("End")
			break
		}

		r := strings.NewReader(string(buffer))

		v, err := parse(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		command := strings.ToUpper(v.array[0].bulk)
		args := v.array[1:]

		handler, ok := handlers[command]
		if !ok {
			// change later
			connection.Write(writeRESP(Value{typ: "string", str: "OK"}))
			continue
		}

		result := handler(args)
		connection.Write(writeRESP(result))
	}
}
