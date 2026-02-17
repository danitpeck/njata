package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "server host")
	port := flag.String("port", "4000", "server port")
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", addr)
	fmt.Println("Type 'quit' to exit\n")

	// Read banner from server
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" {
				fmt.Println(text)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		}
		os.Exit(0)
	}()

	// Read from stdin and send to server
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := stdin.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "quit" {
			fmt.Println("Disconnecting...")
			return
		}

		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
	}
}
