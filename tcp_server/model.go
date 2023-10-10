package tcp_server

import (
	"net"
	"sync"
)

type TCPServer struct {
	buffer_size int32 `default:"1024"`
	isClosed    bool  `default:"false"`
	port        string
	channel     chan []byte

	wg         sync.WaitGroup
	server     net.Listener
	clientList []net.Conn

	handleProcess func(net.Conn, []byte, *TCPServer)
	logger        func([]byte, int)
}
