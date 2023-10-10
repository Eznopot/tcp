package tcp_client

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

func New(port string, handle func([]byte, *TCPClient), optional ...interface{}) *TCPClient {
	var client TCPClient
	client.isClosed = false
	client.channel = make(chan []byte)
	client.handleProcess = handle
	client.port = port

	for _, arg := range optional {
		switch t := arg.(type) {
		case int32:
			client.buffer_size = t
		case func([]byte, int):
			client.logger = t
		default:
			panic("Unknown argument")
		}
	}

	return &client
}

func (c *TCPClient) receivePacket(size int32) ([]byte, error) {
	var received = make([]byte, size)

	err := binary.Read(c.server, binary.LittleEndian, &received)
	if err != nil && !c.isClosed {
		fmt.Println("Error decoding message:", err.Error())
		return nil, err
	} else if c.isClosed {
		return nil, nil
	}
	return received, nil
}

func (c *TCPClient) receiveProcess() {
	for {
		// receive message size of the next packet
		var size int32
		err := binary.Read(c.server, binary.LittleEndian, &size)
		if err != nil && !c.isClosed {
			fmt.Println("Error decoding size:", err.Error())
			return
		} else if c.isClosed {
			return
		}
		if size < 0 {
			println("server disconnected")
			break
		}
		received, err := c.receivePacket(size)
		if err != nil {
			fmt.Println("Error receiving:", err.Error())
			return
		}
		if c.logger != nil {
			c.logger(received, len(received))
		}
		if c.handleProcess != nil {
			c.handleProcess(received, c)
		}
	}
	c.server.Close()
	c.wg.Done()
}

func (c *TCPClient) Send(msg []byte) {
	// send message size of the next packet
	binary.Write(c.server, binary.LittleEndian, int32(len(msg)))
	// send message
	binary.Write(c.server, binary.LittleEndian, msg)
}

func (c *TCPClient) Connect() (*sync.WaitGroup, chan []byte) {
	var err error
	c.server, err = net.Dial("tcp", ":"+c.port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	c.wg.Add(1)
	go c.receiveProcess()
	go func() {
		for {
			msg := <-c.channel
			if c.isClosed {
				break
			}
			c.Send(msg)
		}
	}()
	return &c.wg, c.channel
}

func (c *TCPClient) Close() {
	binary.Write(c.server, binary.LittleEndian, int32(-1))
	c.isClosed = true
	c.server.Close()
	c.channel <- []byte{}
	c.wg.Done()
}
