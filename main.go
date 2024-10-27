package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	EXIT_COMMAND = "exit"
	PORT = ":1729"
)

func main() {
	// Start listening on port 1729
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	// Ensure that listener is closed when the main function exits
	defer listener.Close()

	fmt.Printf("TCP server is listening on port %s ...\n", PORT[1:])

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		
		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

// Function to handle an individual client connection
func handleConnection(conn net.Conn) {

	// anonymous function; which is immediately defined and will be executed when handleConnection returns.
	defer func() {
		conn.Close()
		fmt.Printf("Connection with %v closed\n", conn.RemoteAddr())
	}()

	fmt.Printf("Connected to: %v\n", conn.RemoteAddr())

	// Welcome message with instructions
	welcomeMessage := "Welcome! Type 'exit' to close the connection.\n"

	if _, err := conn.Write([]byte(welcomeMessage)); err != nil {
		log.Printf("Error sending welcom message to %v : %v", conn.RemoteAddr(), err)
		return
	}

	buffer := make([]byte, 2048)

	for {
		// Read data sent by client
		n, err := conn.Read(buffer)
		if err != nil {
			// normal connection closure
			if err != net.ErrClosed {
				log.Printf("Error reading from %v: %v", conn.RemoteAddr(), err)
			}
			return // Exit the goroutine
		}

		// Convert the message to string
		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received from %v: %s\n", conn.RemoteAddr(), message)

		// check if client wants to exit
		if strings.ToLower(message) == EXIT_COMMAND {
			goodByeMessage := "Goodbye! Clossing connection.\n"
			if _, err := conn.Write([]byte(goodByeMessage)); err != nil {
				log.Printf("Error sending goodbye message to %v : %v", conn.RemoteAddr(), err)
			}
			return	// This will trigger the deferred connection close
		}

		// Convert the message to uppercase
		uppercaseMessage := strings.ToUpper(message)

		// Send response
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n%s\r\n", uppercaseMessage)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to %v: %v", conn.RemoteAddr(), err)
			return
		}
	}
}