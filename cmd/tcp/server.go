package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ruhollahh/gopher-cache/internal/command"
	"github.com/ruhollahh/gopher-cache/internal/resp"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	requestReader := bufio.NewReader(conn)
	for {
		request, err := resp.Decode(requestReader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("error reading response: ", err.Error())
		}

		commandHandler := command.NewHandler()
		commandHandler.RegisterHandlers()

		rawResponse := commandHandler.Handle(request)

		encodedResponse, err := resp.Encode(rawResponse)
		if err != nil {
			fmt.Println("error encoding response: ", err.Error())
		}

		_, err = conn.Write(encodedResponse)
		if err != nil {
			fmt.Println("error writing response: ", err.Error())
		}
	}

	err := conn.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
	}
}
