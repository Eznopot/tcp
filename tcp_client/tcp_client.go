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

var isClosed = false

var server net.Conn

var wg sync.WaitGroup

var msgToServer chan []byte

var handleProcess func([]byte)

func receivePacket(size int32) ([]byte, error) {
	var received = make([]byte, size)
	err := binary.Read(server, binary.LittleEndian, &received)
	if err != nil && !isClosed {
		fmt.Println("Error decoding message:", err.Error())
		return nil, err
	} else if isClosed {
		return nil, nil
	}
	return received, nil
}

func receiveProcess() {
	for {
		// receive message size of the next packet
		var size int32
		err := binary.Read(server, binary.LittleEndian, &size)
		if err != nil && !isClosed {
			fmt.Println("Error decoding size:", err.Error())
			return
		} else if isClosed {
			return
		}
		if size < 0 {
			println("server disconnected")
			break
		}
		received, err := receivePacket(size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			return
		}
		handleProcess(received)
	}
	server.Close()
	wg.Done()
}

func Send(msg []byte) {
	// send message size of the next packet
	binary.Write(server, binary.LittleEndian, int32(len(msg)))
	// send message
	binary.Write(server, binary.LittleEndian, msg)
}

func Client(port string, handle func([]byte)) (*sync.WaitGroup, chan []byte) {
	var err error
	handleProcess = handle
	server, err = net.Dial("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	wg.Add(1)
	msgToServer = make(chan []byte)
	go receiveProcess()
	go func() {
		for {
			msg := <-msgToServer
			if isClosed {
				break
			}
			Send(msg)
		}
	}()
	return &wg, msgToServer
}

func Close() {
	binary.Write(server, binary.LittleEndian, int32(-1))
	isClosed = true
	server.Close()
	msgToServer <- []byte{}
	wg.Done()
}
