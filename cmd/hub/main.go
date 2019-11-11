package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"io"
	"os"

	"github.com/martinmarks/chat/pkg/chat"
)

func main() {
	f, err := os.OpenFile("hub.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	cert, err := tls.LoadX509KeyPair("selfsigned.crt", "selfsigned.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":5000", config)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Starting listening for connections...")
	hub := chat.NewHub()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		hub.AcceptConnection(conn)
	}

}
