package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const filePath = "messages.txt"

func main() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("could not open %s: %s\n", filePath, err)
	}

	linesChan := getLinesChannel(f)
	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
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
