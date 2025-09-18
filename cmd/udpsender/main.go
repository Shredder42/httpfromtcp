package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const UDPAddress = "localhost:42069"

func main() {
	UDPAddr, err := net.ResolveUDPAddr("udp", UDPAddress)
	if err != nil {
		log.Fatalf("Could not resolve address %s: %s", UDPAddress, err)
	}

	conn, err := net.DialUDP("udp", nil, UDPAddr)
	if err != nil {
		log.Fatalf("Error dialing UDP: %s", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		line, err := reader.ReadString(('\n'))
		if err != nil {
			log.Fatalf("Error reading string from stdin: %s", err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error writing to the connection: %s", err)
		}

		fmt.Printf("Message sent: %s", line)
	}
}
