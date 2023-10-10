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
	var chanMsg chan []byte
	var wg *sync.WaitGroup

	// close server on ctrl C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tcp_server.Close()
	}()

	//start server
	wg, chanMsg = tcp_server.Server(port, func(client net.Conn, msg []byte) {
		fmt.Printf("Message from client: %s \n", string(msg))

		chanMsg <- []byte("Hello to all from server with chan")          // send msg to all via channel
		tcp_server.SendAll([]byte("Hello to all from server with func")) // send msg to all via function

		//send msg to client
		tcp_server.Send(client, []byte("Hello from server to you")) // send msg to client via function
	})

	//wait for the process to finish
	wg.Wait()
}

func client(port string) {
	fmt.Println("client")

	//connect to server

	// close client on ctrl C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tcp_client.Close()
	}()

	wg, chanMsg := tcp_client.Client(port, func(msg []byte) {
		fmt.Printf("Message from server: %s \n", string(msg))
	})

	//same thing
	chanMsg <- []byte("Hello from client via chan")       // send msg via channel
	tcp_client.Send([]byte("Hello from client via func")) // send msg via function

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
