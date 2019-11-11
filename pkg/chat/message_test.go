package chat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	msg := ParseMessage("1 1 2,3 hello world !")
	assert.Equal(t, "2,3", msg.Receivers)
	assert.Equal(t, "1", msg.Sender)
	assert.Equal(t, TypeUserMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "hello world !"))
}

func TestParseMessageEmpty(t *testing.T) {
	msg := ParseMessage("1 1 2,3")
	assert.Empty(t, msg.Receivers)
	assert.Empty(t, msg.Sender)
	assert.Empty(t, msg.Type)
	assert.Empty(t, msg.Text)
}

func TestUserMessage(t *testing.T) {
	data := "all hey"
	msg := UserMessage(data)
	assert.Equal(t, All, msg.Receivers)
	assert.Equal(t, "user", msg.Sender)
	assert.Equal(t, TypeUserMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "hey"))
	assert.True(t, strings.Contains(msg.Text, ColorCyan))
	assert.True(t, strings.Contains(msg.Text, ColorOff))
}

func TestUserMessageNotAll(t *testing.T) {
	data := "1,2,3,4,5 hey"
	msg := UserMessage(data)
	assert.Equal(t, "1,2,3,4,5", msg.Receivers)
	assert.Equal(t, "user", msg.Sender)
	assert.Equal(t, TypeUserMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "hey"))
	assert.True(t, strings.Contains(msg.Text, ColorCyan))
	assert.True(t, strings.Contains(msg.Text, ColorOff))
}

func TestUserMessageEmpty(t *testing.T) {
	data := "hey"
	msg := UserMessage(data)
	assert.Empty(t, msg.Receivers)
	assert.Empty(t, msg.Sender)
	assert.Empty(t, msg.Type)
	assert.Empty(t, msg.Text)
}

func TestUserListMessage(t *testing.T) {
	data := []string{"one", "two", "three"}
	msg := UserListMessage(data)
	assert.Equal(t, All, msg.Receivers)
	assert.Equal(t, "system", msg.Sender)
	assert.Equal(t, TypeClientListMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "one"))
	assert.True(t, strings.Contains(msg.Text, "two"))
	assert.True(t, strings.Contains(msg.Text, "three"))
	assert.True(t, strings.Contains(msg.Text, ColorBlue))
	assert.True(t, strings.Contains(msg.Text, ColorOff))
}

func TestConnectMessage(t *testing.T) {
	data := "bob"
	msg := ConnectMessage(data)
	assert.Equal(t, All, msg.Receivers)
	assert.Equal(t, "bob", msg.Sender)
	assert.Equal(t, TypeConnectMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "bob connected"))
	assert.True(t, strings.Contains(msg.Text, ColorGreen))
	assert.True(t, strings.Contains(msg.Text, ColorOff))
}

func TestDisconnectMessage(t *testing.T) {
	data := "bob"
	msg := DisconnectMessage(data)
	assert.Equal(t, All, msg.Receivers)
	assert.Equal(t, "bob", msg.Sender)
	assert.Equal(t, TypeDisconnectMessage, msg.Type)
	assert.True(t, strings.Contains(msg.Text, "bob disconnected"))
	assert.True(t, strings.Contains(msg.Text, ColorRed))
	assert.True(t, strings.Contains(msg.Text, ColorOff))
}
