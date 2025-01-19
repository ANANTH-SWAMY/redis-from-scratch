package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on port 6379...")

	initializeRdb()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			resp := bufio.NewReader(connection)
			defer connection.Close()

			for {
				v, err := parse(resp)
				if err != nil {
					if err == io.EOF {
						return
					}
					fmt.Println(err)
				}

				if reflect.DeepEqual(v, Value{}) {
					continue
				}

				command := strings.ToUpper(v.array[0].bulk)
				args := v.array[1:]

				handler, ok := handlers[command]
				if !ok {
					connection.Write(writeRESP(unknownCommand(v.array[0].bulk)))
					continue
				}

				result := handler(args)
				connection.Write(writeRESP(result))
			}
		}()
	}
}
