package chat

import (
	"net"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	_, mockConn := net.Pipe()
	client := NewClient(mockConn)
	assert.Equal(t, client.state, 1)
	assert.NotNil(t, client.outboundMessage)
	assert.NotNil(t, client.inboundMessage)
	assert.NotNil(t, client.disconnect)
	assert.NotNil(t, client.reader)
	assert.NotNil(t, client.writer)
	assert.Equal(t, mockConn, client.conn)
}

func TestClient_ReceiveMessages_Disconnect(t *testing.T) {
	_, mockConn := net.Pipe()
	client := NewClient(mockConn)
	client.disconnect <- true
	assert.Equal(t, client.state, 0)
	mockConn.SetReadDeadline(time.Now())
	_, err := mockConn.Write([]byte{})
	assert.Error(t, err)
}
