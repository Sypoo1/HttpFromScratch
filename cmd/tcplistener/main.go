package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesChan := getLinesChannel(connection)
		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Connection is closed")
	}

}

func getLinesChannel(stream io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer stream.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := stream.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLineContents += parts[i]
				lines <- currentLineContents
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]

		}
	}()
	return lines
}
