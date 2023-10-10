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

var isClosed = false

var handleProcess func(net.Conn, []byte)

var server net.Listener

var msgToAll chan []byte

var wg sync.WaitGroup

func removeClient(client net.Conn) {
	for i, c := range clientList {
		if c == client {
			clientList = append(clientList[:i], clientList[i+1:]...)
			break
		}
	}
}

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
		if err != nil && !isClosed {
			fmt.Println("Error decoding size:", err.Error())
			return
		} else if isClosed {
			return
		}
		if size < 0 {
			println("Client disconnected")
			removeClient(client)
			break
		}
		received, err := receivePacket(client, size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			return
		}
		if handleProcess != nil {
			handleProcess(client, received)
		}
	}
	client.Close()
}

func GetAllClients() []net.Conn {
	return clientList
}

func SendAll(msg []byte) {
	for _, client := range clientList {
		// send message size of the next packet
		binary.Write(client, binary.LittleEndian, int32(len(msg)))
		// send message
		binary.Write(client, binary.LittleEndian, msg)
	}
}

func Send(client net.Conn, msg []byte) {
	// send message size of the next packet
	binary.Write(client, binary.LittleEndian, int32(len(msg)))
	// send message
	binary.Write(client, binary.LittleEndian, msg)
}

func Server(port string, handle func(net.Conn, []byte)) (*sync.WaitGroup, chan []byte) {
	// create tcp server
	var err error
	msgToAll = make(chan []byte)
	server, err = net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	handleProcess = handle

	// listen for connections
	wg.Add(1)
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil && !isClosed {
				fmt.Println("Error accepting:", err.Error())
				os.Exit(1)
			} else if isClosed {
				break
			}
			clientList = append(clientList, conn)
			go receiveProcess(conn)
		}
	}()

	go func() {
		for {
			msg := <-msgToAll
			if isClosed {
				break
			}
			go SendAll(msg)
		}
	}()
	return &wg, msgToAll
}

func Close() {
	isClosed = true
	for _, client := range clientList {
		binary.Write(client, binary.LittleEndian, int32(-1))
		client.Close()
	}
	server.Close()
	msgToAll <- []byte{}
	wg.Done()
}
