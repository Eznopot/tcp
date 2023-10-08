package tcp_client

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

// create tcp client
const BUFFER_SIZE = 1024

var server net.Conn

var handleProcess func([]byte)

func receivePacket(size int32) ([]byte, error) {
	var received = make([]byte, size)
	err := binary.Read(server, binary.LittleEndian, &received)
	if err != nil {
		fmt.Println("Error decoding message:", err.Error())
		return nil, err
	}
	return received, nil
}

func receiveProcess() {
	for {
		// receive message size of the next packet
		var size int32
		err := binary.Read(server, binary.LittleEndian, &size)
		if err != nil {
			fmt.Println("Error decoding size:", err.Error())
			os.Exit(1)
		}
		received, err := receivePacket(size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			os.Exit(1)
		}
		handleProcess(received)
	}
}

func Send(msg string) {
	// send message size of the next packet
	binary.Write(server, binary.LittleEndian, int32(len(msg)))
	// send message
	binary.Write(server, binary.LittleEndian, []byte(msg))
}

func Client(port string, handle func([]byte)) *sync.WaitGroup {
	var err error
	var wg sync.WaitGroup
	handleProcess = handle
	server, err = net.Dial("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	wg.Add(1)
	go receiveProcess()
	return &wg
}
