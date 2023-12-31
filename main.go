package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Eznopot/tcp/tcp_client"
	"github.com/Eznopot/tcp/tcp_server"
)

func server(port string) {
	fmt.Println("server")

	var wg *sync.WaitGroup

	//Create TCP Server
	TCPServer := tcp_server.New(port, func(client net.Conn, msg []byte, TCPServer *tcp_server.TCPServer) {
		fmt.Printf("Message from client: %s \n", string(msg))
		// send msg to all via channel
		TCPServer.SendAll([]byte("Hello to all from server with func")) // send msg to all via function
		//send msg to client
		tcp_server.Send(client, []byte("Hello from server to you")) // send msg to client via function
	})

	// close server on ctrl C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		TCPServer.Close()
	}()

	//start server
	wg = TCPServer.Open()

	//wait for the process to finish
	wg.Wait()
}

func client(port string) {
	fmt.Println("client")

	//Create TCP Client
	TCPClient := tcp_client.New(port, func(msg []byte, TCPClient *tcp_client.TCPClient) {
		fmt.Printf("Message from server: %s \n", string(msg))
	})

	// close client on ctrl C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		TCPClient.Close()
	}()

	//start client
	wg := TCPClient.Connect()

	TCPClient.Send([]byte("Hello from client via func")) // send msg via function
	//wait for the process to finish
	wg.Wait()
}

func main() {
	os.Args = os.Args[1:]
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go [server|client] [port]")
		os.Exit(1)
	}

	if os.Args[0] == "server" {
		server(os.Args[1])
	} else if os.Args[0] == "client" {
		client(os.Args[1])
	}
}
