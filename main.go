package main

import (
	"fmt"
	"io"
	"os"
	"errors"
	"log"
)

const filePath = "messages.txt"

func main() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("could not open %s: %s\n", filePath, err)
	}

	defer f.Close()

	fmt.Printf("Reading data from %s\n", filePath)

	for{
		b := make([]byte, 8, 8)
		n, err := f.Read(b)

		if err != nil{
			if errors.Is(err, io.EOF){
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		str := string(b[:n])
		fmt.Printf("read: %s\n", str)

		fmt.Printf("len=%d\n", len(str))
	}
}
