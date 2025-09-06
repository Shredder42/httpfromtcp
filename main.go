package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	currentLine := ""

	for {
		b := make([]byte, 8, 8)
		n, err := file.Read(b)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("read: %s\n", currentLine)
				os.Exit(0)
			}
			log.Fatal(err)
		}

		readBytes := b[:n]
		parts := bytes.Split(readBytes, []byte{'\n'})

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, string(parts[i]))
			currentLine = ""
		}

		currentLine += string(parts[len(parts)-1])

	}

}
