package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const srv = "127.0.0.1:1234"

func main() {
	//Connect to the server
	conn, err := net.Dial("tcp", srv)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() //Close connection in the end

	//Print proverbs
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("%s \n", message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v '/n'", err)
	}
}
