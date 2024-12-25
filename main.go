package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServerStruct struct {};

// Responsible for creating a server and Open the connection to the client
func(fs *FileServerStruct) StartTheServer() {
	// CREATE THE SERVER
	listen, err := net.Listen("tcp", ":3000");
	// error handling
	if err != nil {
		log.Fatal(err);
	}
	for {
		// Accepting the client request
		conn, err := listen.Accept();
		// error handling
		if err != nil {
			log.Fatal(err);
		}
		// Once a client Connects, the server start a goroutine to handle 
		// that client connection using  ReadTheDataComingFromClient method.
		// allow multiple client to send ther data concurrenctly
		go fs.ReadTheDataComingFromClient(conn);
	}
}

// Responsibility create a client and send the data into chunks to the server.
func ClientThatSendTheData(size int64) error {

	// creating the client and sending a file to the server.
	// dummy file --> slice of bytes [1,2,3,4,4,5,5,,6,6,6,6,,4,,3,3,2,1,1]
	file := make([]byte, size);
	_, err := io.ReadFull(rand.Reader,file);

	if err != nil {
		return err
	}
	// create a client
	conn,err := net.Dial("tcp",":3000");

	// error handling
	if err != nil {
		return err
	}
	// we are sending the data to the server
	binary.Write(conn,binary.LittleEndian, size)
	// file content is sent to the server using the conn.write(file)
	// this is the part where client is streaming the file to the server over the n/w.
	var streamError error;
	var dummyInt int;
	dummyInt ,streamError = conn.Write(file)
	if streamError != nil {
		return err;
	}
	fmt.Println(dummyInt)
	fmt.Println("Have Written the bytes over the network");
	return nil;
}

// Responsibility Recieving the data and logging the data that client is sending..
func (fs *FileServerStruct) ReadTheDataComingFromClient(conn net.Conn) {
	// client has started snding the data...
	// we can create a buffer to to store the client data
	buf := new(bytes.Buffer)

	// Recieve the data in chunks that client is sending
	for {
		var size int64;
		binary.Read(conn, binary.LittleEndian, &size);
		stored, err := io.CopyN(buf,conn,size);
		// error handling
		if err != nil {
			log.Fatal(err);
		}
		file := buf;
		fmt.Println(file.Bytes());
		fmt.Printf("Chunk Recieved %d ",stored);
	}
}

func main() {
	go func () {
		time.Sleep(3 * time.Second)
		ClientThatSendTheData(500000);
	}();
	Server := &FileServerStruct{}
	Server.StartTheServer();
}