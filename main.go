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

	fmt.Println("Listening on port 6379")

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

		fmt.Println(v) // temp

		connection.Write([]byte("+OK\r\n"))
	}
}

