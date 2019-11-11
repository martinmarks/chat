package chat

import (
	"bufio"
	"fmt"
	"net"
)

// Client holds all the information relevant to a client's connection to the hub
type Client struct {
	name            string
	conn            net.Conn
	writer          *bufio.Writer
	reader          *bufio.Reader
	inboundMessage  chan string
	outboundMessage chan string
	disconnect      chan bool
	state           int
	numReceived     int
	numSent         int
}

// NewClient accepts a connection and produces a Client. Go routines are spawned to manage inbound/outbound messages.
func NewClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client := &Client{
		state:           1,
		name:            "user",
		outboundMessage: make(chan string),
		inboundMessage:  make(chan string),
		disconnect:      make(chan bool),
		conn:            conn,
		writer:          writer,
		reader:          reader,
	}

	go client.ReceiveMessages()
	go client.SendMessages()

	return client
}

// ReceiveMessages watches for outbound messages sent from the hub.
func (client *Client) ReceiveMessages() {
	for {
		select {
		case <-client.disconnect:
			if client.conn != nil {
				err := client.conn.Close()
				if err != nil {
					fmt.Printf("error closing client %s\n", err.Error())
				}
			}
			client.state = 0
			break
		default:
			msg := <-client.outboundMessage
			client.numReceived++
			_, err := client.writer.WriteString(msg)
			if err != nil {
				fmt.Printf("error writing message %s\n", err.Error())
			}
			err = client.writer.Flush()
			if err != nil {
				fmt.Printf("error flushing client %s\n", err.Error())
			}
		}
	}
}

// SendMessages sends messages produced from this Client out to the hub.
func (client *Client) SendMessages() {
	for {
		client.numSent++
		msg, err := client.reader.ReadString('\n')
		if err != nil {
			client.inboundMessage <- DisconnectMessage(client.name).ToString()
			client.state = 0
			client.disconnect <- true
			if client.conn != nil {
				err := client.conn.Close()
				if err != nil {
					fmt.Printf("error closing client %s\n", err.Error())
				}
			}
			break
		}

		message := ParseMessage(msg)
		switch message.Type {
		case TypeConnectMessage:
			client.name = message.Sender
			client.inboundMessage <- message.ToString()
		case TypeUserMessage:
			message.Sender = client.name
			client.inboundMessage <- message.ToString()
		default:
			client.inboundMessage <- message.ToString()
		}
	}
}
