package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Eznopot/tcp/tcp_client"
	"github.com/Eznopot/tcp/tcp_server"
)

func server(port string) {
	fmt.Println("server")
	wg := tcp_server.Server(port, func(client net.Conn, msg []byte) {
		fmt.Printf("Message from client: %s \n", string(msg))
		tcp_server.SendAll(string(msg))
	})
	wg.Wait()
}

func client(port string) {
	fmt.Println("client")
	wg := tcp_client.Client(port, func(msg []byte) {
		fmt.Printf("Message from server: %s \n", string(msg))
		tcp_client.Send("Hello from client")
	})
	tcp_client.Send("Hello from client")
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
