package cmd

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	server "github.com/AnhTTx13/tcp-chat-server/internal"
)

var (
	Port int
)

func init() {
	flag.IntVar(&Port, "port", 8080, "Specify port number")
	flag.Parse()
}

func Execute() {
	var addr string = fmt.Sprintf(":%d", Port)

	fmt.Println("Secure TCP socket chat-server using TLS connection.")
	cert, err := server.LoadCerts()
	if err != nil {
		log.Fatal("Cannot load certificates: ", err.Error())
	}
	fmt.Println("Certificates loaded")
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", addr, &config)

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

	fmt.Println("Server listening on port:", Port)
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
