package rcon

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	packetIDAuthFailed = -1
	payloadMaxSize     = 2048
	payloadTypeAuth    = 3
	payloadTypeCommand = 2
)

// Client represents an RCON client
type Client struct {
	conn     net.Conn
	password string
}

// NewClient initializes a new RCON client
func NewClient(host string, port int16, pass string) (*Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:     conn,
		password: pass,
	}

	return client, nil
}

// Authenticate sends authentication requests to the RCON server
func (c *Client) Authenticate() error {
	payload := newPayload(payloadTypeAuth, c.password)

	if _, err := c.sendPayload(payload); err != nil {
		return err
	}

	return nil
}

// ExecCommand sends a command to the RCON server
func (c *Client) ExecCommand(command string) (string, error) {
	payload := newPayload(payloadTypeCommand, command)

	response, err := c.sendPayload(payload)
	if err != nil {
		return "", err
	}

	// Trim null bytes
	// response.Body = bytes.Trim(response.Body, "\x00")

	return strings.TrimSpace(string(response.Body)), nil
}

// Reconnect tries to reconnect to an RCON server after being disconnected
func (c *Client) Reconnect() error {
	conn, err := net.DialTimeout("tcp", c.conn.RemoteAddr().String(), 10*time.Second)
	if err != nil {
		return err
	}

	c.conn = conn

	err = c.Authenticate()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) sendPayload(request *payload) (*payload, error) {
	packet, err := buildPacketFromPayload(request)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(packet)
	if err != nil {
		return nil, err
	}

	response, err := buildPayloadFromPacket(c.conn)
	if err != nil {
		return nil, err
	}

	if response.ID == packetIDAuthFailed {
		return nil, fmt.Errorf("Authentication failed")
	}

	return response, nil
}
