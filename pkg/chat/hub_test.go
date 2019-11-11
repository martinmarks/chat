package chat

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHub_NewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub.clientMap)
	assert.NotNil(t, hub.connection)
	assert.NotNil(t, hub.outboundMessage)
	assert.Equal(t, 0, hub.numClients)
}

func TestHub_Join(t *testing.T) {
	hub := NewHub()

	_, mockConn := net.Pipe()

	hub.Join(mockConn)
	assert.Equal(t, 1, hub.numClients)
	hub.Join(mockConn)
	assert.Equal(t, 2, hub.numClients)

	_, ok := hub.clientMap["1"]
	assert.True(t, ok)
}

func TestHub_BroadcastUserList(t *testing.T) {
	hub := NewHub()
	_, mockConn := net.Pipe()
	hub.Join(mockConn)
	hub.Join(mockConn)

	client, _ := hub.clientMap["1"]
	client.state = 0
	hub.BroadcastUserList()
	client, ok := hub.clientMap["1"]
	assert.False(t, ok)
}
