package chat

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// Hub manages connected Clients and passes Messages between them
type Hub struct {
	numClients      int
	clientMap       map[string]*Client
	connection      chan net.Conn
	outboundMessage chan string
}

// NewHub initalizes and starts a new hub
func NewHub() *Hub {
	hub := &Hub{
		numClients:      0,
		clientMap:       make(map[string]*Client, 0),
		connection:      make(chan net.Conn),
		outboundMessage: make(chan string),
	}

	hub.Start()

	return hub
}

// Start creates a goroutine that accepts new connections and watches for messages to send.
func (hub *Hub) Start() {
	go func() {
		for {
			select {
			case conn := <-hub.connection:
				hub.Join(conn)
			case msg := <-hub.outboundMessage:
				hub.BroadcastMessage(msg)
			}
		}
	}()
}

// AcceptConnection takes a new connection and passes it to our channel for handling
func (hub *Hub) AcceptConnection(conn net.Conn) {
	hub.connection <- conn
}

// Join creates new client and starts listening for client messages.
func (hub *Hub) Join(conn net.Conn) {
	client := NewClient(conn)
	hub.numClients++
	clientID := strconv.Itoa(hub.numClients)
	hub.clientMap[clientID] = client
	go func() {
		for {
			hub.outboundMessage <- <-client.inboundMessage
		}
	}()
}

// BroadcastUserList sends a message to every connected client with an up-to-date list of connected users
func (hub *Hub) BroadcastUserList() {
	var connectedClients []string
	for id, client := range hub.clientMap {
		if client.state == 0 {
			delete(hub.clientMap, id)
		} else {
			connectedClients = append(connectedClients, fmt.Sprintf("%s-%s", id, client.name))
		}
	}

	for _, client := range hub.clientMap {
		client.outboundMessage <- UserListMessage(connectedClients).ToString()
	}
}

// BroadcastMessage sends a message to the appropriate recipients
func (hub *Hub) BroadcastMessage(data string) {
	hub.BroadcastUserList()
	msg := ParseMessage(data)
	fmt.Println(msg.DisplayString())
	log.Println(msg.DisplayString())

	if msg.Receivers == All {
		for _, client := range hub.clientMap {
			client.outboundMessage <- msg.ToString()
		}

	} else {
		for _, clientID := range strings.Split(msg.Receivers, ",") {
			client, ok := hub.clientMap[clientID]
			if ok {
				client.outboundMessage <- msg.ToString()
			}
		}
	}
}
