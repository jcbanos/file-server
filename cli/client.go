package main

import (
	"net"
	"fmt"
    "os"
    "log"
    "strconv"
    "strings"
    "io"
    "path/filepath"
)

// Constants for client
const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
	BUFFER = 1024
)

func main() {

    // Get client sys args
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

    defer conn.Close()
    
    // Handle the action with appropriate function
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

    // create a directory for client file
    os.Mkdir("client" + "-" + clientIdStr, os.ModePerm)

    fmt.Printf("I have been assigned client id %d and created a directory\n", clientId)
    fmt.Printf("Waiting for files!\n")

    // Wait and recieve file name
    for {
        bufferFileName := make([]byte, 64)
        conn.Read(bufferFileName)
        fileName := strings.Trim(string(bufferFileName), ":")
        println()
        fmt.Printf("Recieving %s from channel %s\n", fileName, arguments[3])
    
        // Create file
        newFile, err := os.Create(filepath.Join("client" + "-" + clientIdStr, filepath.Base(fileName)))
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
    
        defer newFile.Close()
    
        // Wait and recieve file size
        bufferFileSize := make([]byte, 10)
        conn.Read(bufferFileSize)
    
        fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
        fmt.Printf("The file size is: %d\n", fileSize)
    
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
    }
}

func handleSend(conn net.Conn, arguments []string){

    // Fill and send arguments to server
    sendArguments := arguments[1] + " " + arguments[2] + " " + arguments[4]
    sendArguments = fillString(sendArguments, BUFFER)
    fmt.Fprintf(conn, sendArguments)

    // Open file to send
    file, err := os.Open(arguments[2])
	if err != nil {
		fmt.Println(err)
		return
	}


    // Get file info to get file size and send to server
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

    fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
    conn.Write([]byte(fileSize))

    // Create buffer to send file
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

// Function to fill strings with : up to toLength bytes
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
    -How can client read without making a blocking call so that the client can do other things
    while waiting for files

    -What to do when sending files that do not exist
*/