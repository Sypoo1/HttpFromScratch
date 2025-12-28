package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/requests"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection is accepted")
		r, err := requests.RequestFromReader(connection)
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println("Connection is closed")
	}

}
