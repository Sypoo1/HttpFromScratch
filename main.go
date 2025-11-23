package main

import (
	"fmt"
	"io"
	"os"
	"errors"
	"log"
	"strings"
)

const filePath = "messages.txt"

func main() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("could not open %s: %s\n", filePath, err)
	}

	defer f.Close()

	fmt.Printf("Reading data from %s\n", filePath)

	currentLineContents := ""
	for{
		buffer := make([]byte, 8, 8)
		n, err := f.Read(buffer)
		if err != nil{
			if currentLineContents != "" {
				fmt.Printf("read: %s\n", currentLineContents)
				currentLineContents = ""
			}
			if errors.Is(err, io.EOF){
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		str := string(buffer[:n])
		parts := strings.Split(str, "\n")
		for i := 0; i < len(parts) - 1; i++ {
			currentLineContents += parts[i]
			fmt.Printf("read: %s\n", currentLineContents)
			currentLineContents = ""
		}
		currentLineContents += parts[len(parts) - 1]

	}
}
