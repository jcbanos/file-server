package main

import (
	"net"
	"fmt"
    "os"
    "log"
    "strconv"
    "strings"
    "io"
)

// Constants for client
const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
	BUFFER = 1024
)

func main() {

    arguments := os.Args

	if len(arguments) != 4  && len(arguments) != 5 {
		println("Please provide  3 or 4 arguments. Example: ./client recieve -channel 1")
		os.Exit(1)
	}

    // Connect to the server
    conn, err := net.Dial("tcp", HOST+":"+PORT)
    if err != nil {
        fmt.Println("Error accepting: ", err.Error())
        return
    }
    

    if arguments[1] == "recieve" {
        handleRecieve(conn, arguments)
    } else if arguments[1] == "send" {
        handleSend(conn, arguments)
    } else {
        println("Incorrect first argument. Has to be send or recieve")
        os.Exit(1)
    }
 
}

func handleRecieve( conn net.Conn, arguments []string) {
    recvBuf := make([]byte, BUFFER)
    fmt.Fprintf(conn, arguments[1] + " " + arguments[2] + " " + arguments[3])

    // Get client ID
    _, err := conn.Read(recvBuf[:])
    if err != nil {
        fmt.Println("Error reading from server connection")
        os.Exit(1)
    }

    clientIdStr := string(recvBuf)
    clientIdStr = strings.Replace(clientIdStr, "\x00", "", -1)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    clientId, err := strconv.Atoi(clientIdStr)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    fmt.Printf(" I have been assigned client id %d\n", clientId)

    // Wait for file
    _, err = conn.Read(recvBuf[:])
    if err != nil {
        fmt.Println("Error reading from server connection")
        os.Exit(1)
    }
}

func handleSend(conn net.Conn, arguments []string){
    sendArguments := arguments[1] + " " + arguments[2] + " " + arguments[4]
    sendArguments = fillString(sendArguments, BUFFER)
    fmt.Fprintf(conn, sendArguments)

    file, err := os.Open(arguments[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

    fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
    conn.Write([]byte(fileSize))

    sendBuffer := make([]byte, BUFFER)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
        
	}
    fmt.Println("File has been sent, closing connection!")
    return
}

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
