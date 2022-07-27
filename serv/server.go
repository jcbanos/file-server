package main

import (
	"os"
	"net"
	"fmt"
	"strings"
	"strconv"
	"log"
	"io"
	"path/filepath"
)

// Constants for the server
const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
	BUFFER = 1024
)

// Struct for a server. The only attribute is clients a map of slices for io.Writers
type server struct {
    clients map[int][]net.Conn
	clientCount int
}

func (s *server) addClient(c net.Conn, channel int) {
    s.clients[channel] = append(s.clients[channel], c)
}

func (s *server) sendFileToSubscribedClients(fileName string, channel int) {
	clientsToSend := s.clients[channel]
	fmt.Printf("Clients to send to: %d\n", clientsToSend)
	for _, client := range clientsToSend{
		fmt.Printf("%p\n", client )
	}

}

func (s *server) handleClient(c net.Conn) {
	recvBuf := make([]byte, BUFFER)
	_, err := c.Read(recvBuf[:])
	if err != nil {
		fmt.Println("Error reading from client connection")
		os.Exit(1)
	}

	
	clientArgs := strings.Trim(string(recvBuf[:]), ":")
	sliceArgs := strings.Split(clientArgs, " ")
	if sliceArgs[0] == "recieve" {
		channel, err := strconv.Atoi(strings.Replace(sliceArgs[2], "\x00", "", -1))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		strId:= strconv.Itoa(s.clientCount)
		fmt.Fprintf(c, strId)
		s.clientCount += 1
		s.addClient(c, channel)
	} else if sliceArgs[0] == "send" {
		fmt.Printf("%v\n", sliceArgs)
		channel, err := strconv.Atoi(strings.Replace(sliceArgs[2], "\x00", "", -1))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		strId:= strconv.Itoa(s.clientCount)
		fmt.Fprintf(c, strId)
		s.clientCount += 1
		fmt.Printf("Client #%d is sending file %s to channel %d\n", s.clientCount, sliceArgs[1], channel)
		bufferFileSize := make([]byte, 10)

		c.Read(bufferFileSize)
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
		fmt.Printf("The file size is: %d\n", fileSize)

		newFile, err := os.Create(filepath.Join("files", filepath.Base(sliceArgs[1])))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer newFile.Close()
		var receivedBytes int64

		for {
			if (fileSize - receivedBytes) < BUFFER {
				io.CopyN(newFile, c, (fileSize - receivedBytes))
				c.Read(make([]byte, (receivedBytes+BUFFER)-fileSize))
				break
			}
			io.CopyN(newFile, c, BUFFER)
			receivedBytes += BUFFER
		}
		fmt.Println("Received file completely!")
		s.sendFileToSubscribedClients(sliceArgs[1], channel)
	}
}	

func main(){

	// Check the args to see if client passes start
	arguments := os.Args
	if len(arguments) != 2 {
		println("Please provide 1 argument. Example: ./server start")
		os.Exit(1)
	}

	// If argument 1 is start, then start server otherwise shutdown
	if arguments[1] != "start" {
		println("Shutting down server! Goodbye!")
		return
	} else {
		fmt.Println("Launching server...")
	}

	// Create server struct with maps of client connections to keep track of client's subscriptions
	srv := &server{}
	srv.clients = make(map[int][]net.Conn)
	srv.clientCount = 0

	// Start to listen for incoming connections
	ln, _ := net.Listen("tcp", HOST+":"+PORT)

	for {
        // Accept incoming connections
        conn, err := ln.Accept()

        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            continue
        }

		// Handle connection with handleClient Function
		go srv.handleClient(conn)

    }
}


