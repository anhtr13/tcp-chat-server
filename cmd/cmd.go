package cmd

import (
	"flag"
	"fmt"
	"log"
	"net"

	server "github.com/AnhBigBrother/tcp-chat-server/internal"
)

var PORT int

func init() {
	flag.IntVar(&PORT, "port", 8080, "Specify port number")
	flag.Parse()
}

func Execute() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer listener.Close()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil { // Check for IPv4
				fmt.Printf("IPv4 Address: %s\n", ipnet.IP.String())
			} else if ipnet.IP.To16() != nil { // Check for IPv6
				fmt.Printf("IPv6 Address: %s\n", ipnet.IP.String())
			}
		}
	}

	fmt.Println("Socket server listening on port:", PORT)
	s := server.NewServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.HandleConnection(conn)
	}
}
