package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/Eznopot/tcp/tcp_client"
	"github.com/Eznopot/tcp/tcp_server"
)

func server(port string) {
	fmt.Println("server")
	var chanMsg chan string
	var wg *sync.WaitGroup

	//start server
	wg, chanMsg = tcp_server.Server(port, func(client net.Conn, msg []byte) {
		fmt.Printf("Message from client: %s \n", string(msg))

		chanMsg <- "Hello to all from server with chan"          // send msg to all via channel
		tcp_server.SendAll("Hello to all from server with func") // send msg to all via function

		//send msg to client
		tcp_server.Send(client, "Hello from server to you") // send msg to client via function
	})
	wg.Wait()
}

func client(port string) {
	fmt.Println("client")

	//connect to server
	wg, chanMsg := tcp_client.Client(port, func(msg []byte) {
		fmt.Printf("Message from server: %s \n", string(msg))
	})

	//same thing
	chanMsg <- "Hello from client via chan"       // send msg via channel
	tcp_client.Send("Hello from client via func") // send msg via function
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
