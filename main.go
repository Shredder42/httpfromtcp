package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {

	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()

		currentLine := ""

		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}
				if err == io.EOF {
					return
				}
				fmt.Printf("error: %s\n", err)
			}

			readBytes := b[:n]
			parts := bytes.Split(readBytes, []byte{'\n'})

			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLine, string(parts[i]))
				currentLine = ""
			}

			currentLine += string(parts[len(parts)-1])

		}

	}()

	return lines

}

func main() {

	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)

	linesChannel := getLinesChannel(f)

	for line := range linesChannel {
		fmt.Printf("read: %s\n", line)
	}

}
