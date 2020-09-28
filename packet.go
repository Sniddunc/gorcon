package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func buildPacketFromPayload(payload *payload) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// Write payload data into buffer using LittleEndian as specified in the
	// Source RCON specification.
	binary.Write(buffer, binary.LittleEndian, payload.getSize())
	binary.Write(buffer, binary.LittleEndian, payload.ID)
	binary.Write(buffer, binary.LittleEndian, payload.Type)
	binary.Write(buffer, binary.LittleEndian, payload.Body)
	binary.Write(buffer, binary.LittleEndian, [2]byte{}) // write null bytes

	if buffer.Len() >= payloadMaxSize {
		return nil, fmt.Errorf("Payload too large. Max size: %d", payloadMaxSize)
	}

	return buffer.Bytes(), nil
}

func buildPayloadFromPacket(reader io.Reader) (*payload, error) {
	var packetSize int32
	var packetID int32
	var packetType int32

	// Read header bytes
	err := binary.Read(reader, binary.LittleEndian, &packetSize)
	if err != nil {
		return nil, fmt.Errorf("could not read packet bytes")
	}

	err = binary.Read(reader, binary.LittleEndian, &packetID)
	if err != nil {
		return nil, fmt.Errorf("could not read packet bytes")
	}

	err = binary.Read(reader, binary.LittleEndian, &packetType)
	if err != nil {
		return nil, fmt.Errorf("could not read packet bytes")
	}

	// Create byte slice to read the body into
	packetBody := make([]byte, packetSize-(payloadIDBytes+payloadTypeBytes))

	_, err = io.ReadFull(reader, packetBody)
	if err != nil {
		return nil, err
	}

	payload := &payload{
		ID:   packetID,
		Type: packetType,
		Body: packetBody,
	}

	return payload, nil
}
