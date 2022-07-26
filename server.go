package main

import (
	"os"
	"net"
	"fmt"
	"log"
	"time"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
)

func main(){

	arguments := os.Args
	if len(arguments) != 2 {
		println("Please provide 1 argument. Example: ./server start")
		os.Exit(1)
	}

	if arguments[1] != "start" {
		println("Shutting down server! Goodbye!")
		return
	} else {
		println("Starting server!")
	}

	listen, err := net.Listen(TYPE, HOST+":"+PORT)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		} 
		go handleIncomingRequest(conn)
	}
}

func handleIncomingRequest(conn net.Conn) {
    // store incoming data
    buffer := make([]byte, 1024)
    _, err := conn.Read(buffer)
    if err != nil {
        log.Fatal(err)
    }
    // respond
    time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
    conn.Write([]byte("Hi back!"))
    conn.Write([]byte(time))
    // close conn
    conn.Close()
}
