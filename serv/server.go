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

// Struct for a server. The attributes are:
// clients: a map of slices for net.Conns
// clientCount: integer to keep count of clients
type server struct {
    clients map[int][]net.Conn
	clientCount int
}

// instance method to add a client to the clients map.
func (s *server) addClient(c net.Conn, channel int) {
    s.clients[channel] = append(s.clients[channel], c)
}


// Function that will send file to cleints subscribed to a channel
func (s *server) sendFileToSubscribedClients(fileName string, channel int) {

	// Retrieve list of clients whom I will send file to
	clientsToSend := s.clients[channel]

	for _, clientConn := range clientsToSend {
		// Open file to be sent
		file, err := os.Open(filepath.Join("files", filepath.Base(fileName)))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get file info to get file size and file name. Then send to clients
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}

		fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
		fileNameSend := fillString(fileInfo.Name(), 64)
		// File name to client
		clientConn.Write([]byte(fileNameSend))
		// Filesize to client
		clientConn.Write([]byte(fileSize))
		// Create buffer to send file
		sendBuffer := make([]byte, BUFFER)

		// Send file
		for {
			_, err := file.Read(sendBuffer)
			if err == io.EOF {
				break
			}
			clientConn.Write(sendBuffer)
		}
		file.Close()
	}

	fmt.Println(fileName + " has been sent to clients in channel " + strconv.Itoa(channel) + "!")
	println()
}


// Function that handles a client on initial connection
func (s *server) handleClient(conn net.Conn) {

	// Create a buffer
	recvBuf := make([]byte, BUFFER)

	// read from the connection to get client args
	_, err := conn.Read(recvBuf[:])
	if err != nil {
		fmt.Println("Error reading from client connection")
		os.Exit(1)
	}

	// Process clients args by removing delimitir : and splitting the string
	clientArgs := strings.Trim(string(recvBuf[:]), ":")
	sliceArgs := strings.Split(clientArgs, " ")

	// if client is doing recieve action then handle recieve
	if sliceArgs[0] == "recieve" {
		s.handleRecieveClient(conn, sliceArgs)
	} else if sliceArgs[0] == "send" {
		s.handleSendClient(conn, sliceArgs)
	}
}

// Method to handle clients with recieve keyword
func (s *server) handleRecieveClient(conn net.Conn, sliceArgs []string){

	// Get channel number from client args
	channel, err := strconv.Atoi(strings.Replace(sliceArgs[2], "\x00", "", -1))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Get a client Id from clientcount
	strId:= strconv.Itoa(s.clientCount)

	// Send the client id to client
	fmt.Fprintf(conn, strId)

	// Increment client count and add client connection/channel to the map
	s.clientCount += 1
	s.addClient(conn, channel)

	fmt.Printf("Client with id %s is now subscribed to channel %s\n", strId, sliceArgs[2])
	println()
}

func (s *server) handleSendClient(conn net.Conn, sliceArgs []string){

	// Get the channel number
	channel, err := strconv.Atoi(strings.Replace(sliceArgs[2], "\x00", "", -1))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//Get an id for the client send it to client and increment client count
	strId:= strconv.Itoa(s.clientCount)
	fmt.Fprintf(conn, strId)
	s.clientCount += 1
	fmt.Printf("Client #%d is sending file %s to channel %d\n", s.clientCount, sliceArgs[1], channel)

	// Create buffer to get file size and read connection for file size
	bufferFileSize := make([]byte, 10)
	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Printf("The file size is: %d\n", fileSize)

	// Create a file to download the clients file
	newFile, err := os.Create(filepath.Join("files", filepath.Base(sliceArgs[1])))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer newFile.Close()

	// Create variable to keep track of how many bytes have been recieved
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFER {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFER)-fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFER)
		receivedBytes += BUFFER
	}
	fmt.Println("Received file completely!")

	// After downloading file, send file to clients subscribed to channel
	s.sendFileToSubscribedClients(sliceArgs[1], channel)
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
		println()
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

		defer conn.Close()

		// Handle connection with handleClient Function
		go srv.handleClient(conn)

    }
}

// Function to fill strings with : before sendind the message
func fillString(returnString string, toLength int) string {
	for {
		lengtString := len(returnString)
		if lengtString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}

/* Things I would like to improve and discuss with teamates:

	-Use a mutex with the server struct since there could be race conditions

	-Instead of writing a concurrent server, use worker threads to make server run in parallel

	-Handle closed connections
*/

