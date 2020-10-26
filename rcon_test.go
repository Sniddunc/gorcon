package rcon

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var (
	host     string
	port     int16
	password string
)

func init() {
	// Load environment variables.
	// This is required because to properly test RCON, we need to be able to create a connection.
	// The .env file should have RCON_HOST, RCON_PORT, and RCON_PASSWORD set to a valid rcon server's details.
	// I'm fully aware that this is not ideal testing since the tests are not independent from external sources
	// but that's a risk I'm willing to take for this project as it's so simple, at least in it's current form.
	//
	// Ideally, we would have interfaces representing the connection streams which we could swap a mock one into
	// in order to run the tests against a fake server.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file for testing. Error: %v", err)
	}

	envHost, exists := os.LookupEnv("RCON_HOST")
	if !exists {
		log.Fatalf("The required env variable RCON_HOST is not set")
	}

	envPort, exists := os.LookupEnv("RCON_PORT")
	if !exists {
		log.Fatalf("The required env variable RCON_PORT is not set")
	}

	envPassword, exists := os.LookupEnv("RCON_PASSWORD")
	if !exists {
		log.Fatalf("The required env variable RCON_PASSWORD is not set")
	}

	parsedPort, err := strconv.Atoi(envPort)
	if err != nil {
		log.Fatalf("The value of env variable RCON_PORT is not a valid integer")
	}

	host = envHost
	port = int16(parsedPort)
	password = envPassword
}

func TestClient(t *testing.T) {
	client, err := NewClient(host, port, password)
	if err != nil {
		t.Errorf("NewClient shouldn't return an error, but it did. Error: %v", err)
	}

	err = client.Authenticate()
	if err != nil {
		t.Errorf("Authentication should've passed, but it failed. Error: %v", err)
	}

	res, err := client.ExecCommand("help")
	if err != nil {
		t.Errorf("Command execution should've succeeded, but it failed. Error: %v", err)
	}

	fmt.Println(res)

	res, err = client.ExecCommand(strings.Repeat("a", payloadMaxSize+1))
	if err == nil {
		t.Errorf("Payload was too big, but was sent to the server anyway")
	}

	client.Reconnect()

	fmt.Println(res)
}

func TestNewPayload(t *testing.T) {
	payload := newPayload(payloadTypeCommand, "help")

	if payload.Type != payloadTypeCommand {
		t.Errorf("")
	}
}

func TestPacket(t *testing.T) {
	// Determined via manual calculation
	const properSize = 14

	payload := newPayload(payloadTypeCommand, "help")

	packet, err := buildPacketFromPayload(payload)
	if err != nil {
		t.Errorf("Could not build packet from payload. Error: %v", err)
	}

	// Manually read bytes out from the packet
	pSize := int32(binary.LittleEndian.Uint32(packet[0:4]))
	pID := int32(binary.LittleEndian.Uint32(packet[4:8]))
	pType := int32(binary.LittleEndian.Uint32(packet[8:12]))
	pBody := packet[12 : len(packet)-2]

	if pSize != properSize {
		t.Errorf("Packet size contained the wrong size")
	}

	if pID != int32(currentPayloadID) {
		t.Errorf("Packet id contained the wrong ID. Expected: %d, actual: %d", currentPayloadID, pID)
	}

	if pType != payloadTypeCommand {
		t.Errorf("Packet type contained the wrong type. Expected: %d, actual: %d", payloadTypeCommand, pType)
	}

	if string(pBody) != string(payload.Body) {
		t.Errorf("Packet body contained incorrect data. Expected: %s, actual: %s", payload.Body, pBody)
	}
}
