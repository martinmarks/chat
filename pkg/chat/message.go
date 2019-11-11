package chat

import (
	"fmt"
	"strings"
	"time"
)

// Message type constants
const (
	TypeUserMessage       = "1"
	TypeConnectMessage    = "2"
	TypeDisconnectMessage = "3"
	TypeClientListMessage = "4"

	All = "all"

	ColorOff   = "\033[0m"
	ColorGreen = "\x1b[0;32m"
	ColorRed   = "\x1b[0;31m"
	ColorBlue  = "\x1b[0;34m"
	ColorCyan  = "\x1b[0;36m"
)

// Message holds information about all communications sent between Hub and Clients
type Message struct {
	Text      string
	Type      string
	Sender    string
	Receivers string
}

// UserMessage is a message sent from one User/Client to other(s)
func UserMessage(data string) Message {
	tokens := strings.Split(data, " ")

	if len(tokens) < 2 {
		return Message{}
	}

	return Message{
		Text:      fmt.Sprintf("%s%s%s", ColorCyan, strings.Join(tokens[1:], " "), ColorOff),
		Type:      TypeUserMessage,
		Receivers: tokens[0],
		Sender:    "user",
	}
}

// DisconnectMessage gives a message used when a  user disconnects from the hub
func DisconnectMessage(name string) Message {
	return Message{
		Text:      fmt.Sprintf("%s %s disconnected %s", ColorRed, name, ColorOff),
		Type:      TypeDisconnectMessage,
		Receivers: All,
		Sender:    name,
	}
}

// ConnectMessage gives a message used when a new user connects to the hub
func ConnectMessage(name string) Message {
	return Message{
		Text:      fmt.Sprintf("%s %s connected %s", ColorGreen, name, ColorOff),
		Type:      TypeConnectMessage,
		Receivers: All,
		Sender:    name,
	}
}

// UserListMessage gives a message used when updating the list of currently online users
func UserListMessage(users []string) Message {
	userList := ""
	for _, user := range users {
		userList += fmt.Sprintf("%s%s%s ", ColorBlue, user, ColorOff)
	}

	return Message{
		Text:      userList,
		Type:      TypeClientListMessage,
		Receivers: All,
		Sender:    "system",
	}
}

// DisplayString gives a nicely formatted version of this message ready to displayed on the UI
func (m Message) DisplayString() string {
	currentTime := time.Now().Format("15:04:05")
	return fmt.Sprintf("[%s] [From: %s] [To: %s] %s\n", currentTime, m.Sender, m.Receivers, m.Text)
}

// ToString gives a simple representation of this message. Easy for reading/writing over the wire.
func (m Message) ToString() string {
	return fmt.Sprintf("%s %s %s %s\n", m.Type, m.Sender, m.Receivers, m.Text)
}

// ParseMessage takes a message in string form and returns a Message struct
func ParseMessage(msg string) Message {
	tokens := strings.Fields(msg)

	if len(tokens) < 4 {
		return Message{}
	}

	return Message{Type: tokens[0], Sender: tokens[1], Receivers: tokens[2], Text: strings.Join(tokens[3:], " ")}
}
