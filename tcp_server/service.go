package tcp_server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

func New(port string, handle func(net.Conn, []byte, *TCPServer), optional ...interface{}) *TCPServer {
	var server TCPServer
	server.isClosed = false
	server.handleProcess = handle
	server.port = port

	for _, arg := range optional {
		switch t := arg.(type) {
		case int32:
			server.buffer_size = t
		case func([]byte, int):
			server.logger = t
		default:
			panic("Unknown argument")
		}
	}

	return &server
}

func (s *TCPServer) removeClient(client net.Conn) {
	for i, c := range s.clientList {
		if c == client {
			s.clientList = append(s.clientList[:i], s.clientList[i+1:]...)
			break
		}
	}
}

func (s *TCPServer) receivePacket(client net.Conn, size int32) ([]byte, error) {
	var received = make([]byte, size)
	err := binary.Read(client, binary.LittleEndian, &received)
	if err != nil {
		fmt.Println("Error decoding message:", err.Error())
		return nil, err
	}
	return received, nil
}

func (s *TCPServer) receiveProcess(client net.Conn) {
	for {
		// receive message size of the next packet
		var size int32
		err := binary.Read(client, binary.LittleEndian, &size)
		if err != nil && !s.isClosed {
			fmt.Println("Error decoding size:", err.Error())
			return
		} else if s.isClosed {
			return
		}
		if size < 0 {
			println("Client disconnected")
			s.removeClient(client)
			break
		}
		received, err := s.receivePacket(client, size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			return
		}
		if s.logger != nil {
			s.logger(received, len(received))
		}
		if s.handleProcess != nil {
			s.handleProcess(client, received, s)
		}
	}
	client.Close()
}

func (s *TCPServer) GetAllClients() []net.Conn {
	return s.clientList
}

func (s *TCPServer) SendAll(msg []byte) {
	for _, client := range s.clientList {
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

func (s *TCPServer) Open() *sync.WaitGroup {
	// create tcp server
	var err error
	s.server, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// listen for connections
	s.wg.Add(1)
	go func() {
		for {
			conn, err := s.server.Accept()
			if err != nil && !s.isClosed {
				fmt.Println("Error accepting:", err.Error())
				os.Exit(1)
			} else if s.isClosed {
				break
			}
			s.clientList = append(s.clientList, conn)
			go s.receiveProcess(conn)
		}
	}()

	return &s.wg
}

func (s *TCPServer) Close() {
	s.isClosed = true
	for _, client := range s.clientList {
		binary.Write(client, binary.LittleEndian, int32(-1))
		client.Close()
	}
	s.server.Close()
	s.wg.Done()
}
