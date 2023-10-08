package tcp_server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

const BUFFER_SIZE = 1024

var clientList []net.Conn

var handleProcess func(net.Conn, []byte)

func receivePacket(client net.Conn, size int32) ([]byte, error) {
	var received = make([]byte, size)
	err := binary.Read(client, binary.LittleEndian, &received)
	if err != nil {
		fmt.Println("Error decoding message:", err.Error())
		return nil, err
	}
	return received, nil
}

func receiveProcess(client net.Conn) {
	for {
		// receive message size of the next packet
		var size int32
		err := binary.Read(client, binary.LittleEndian, &size)
		if err != nil {
			fmt.Println("Error decoding size:", err.Error())
			os.Exit(1)
		}
		received, err := receivePacket(client, size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			os.Exit(1)
		}
		if handleProcess != nil {
			handleProcess(client, received)
		}
	}
}

func SendAll(msg string) {
	for _, client := range clientList {
		// send message size of the next packet
		binary.Write(client, binary.LittleEndian, int32(len(msg)))
		// send message
		binary.Write(client, binary.LittleEndian, []byte(msg))
	}
}

func Server(port string, handle func(net.Conn, []byte)) *sync.WaitGroup {
	// create tcp server
	var wg sync.WaitGroup
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	handleProcess = handle

	// listen for connections
	wg.Add(1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err.Error())
				os.Exit(1)
			}
			clientList = append(clientList, conn)
			go receiveProcess(conn)
		}
	}()
	return &wg
}
