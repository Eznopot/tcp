package tcp_client

import (
	"net"
	"sync"
)

type TCPClient struct {
	// create tcp client
	buffer_size   int32 `default:"1024"`
	isClosed      bool  `default:"false"`
	server        net.Conn
	wg            sync.WaitGroup
	handleProcess func([]byte, *TCPClient)
	logger        func([]byte, int)
	port          string
}
