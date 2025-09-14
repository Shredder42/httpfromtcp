package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

const port = ":42069"

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

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not create listener: %s\n", err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not establish connection: %s", err)
		}
		fmt.Println("Connection established from", conn.RemoteAddr())

		linesChannel := getLinesChannel(conn)

		// don't necessarily need a go routine here
		go func() {
			for {
				// this manually checks for the channel closure
				line, ok := <-linesChannel
				if !ok {
					conn.Close()
					fmt.Printf("connection to %s closed\n", conn.RemoteAddr())
					return
				}
				fmt.Println(line)
			}
		}()
	}
}
