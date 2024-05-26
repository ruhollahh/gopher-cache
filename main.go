package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("starting server")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("error accepting connection: ", err.Error())
		os.Exit(1)
	}

	commandBuffer := make([]byte, 1024)
	for {
		if _, err = conn.Read(commandBuffer); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("error reading response: ", err.Error())
		}

		if _, err = conn.Write([]byte("+PONG\r\n")); err != nil {
			fmt.Println("error writing response: ", err.Error())
		}
	}

	err = conn.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
	}
}
