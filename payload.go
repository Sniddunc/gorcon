package rcon

const (
	payloadIDBytes   = 4
	payloadTypeBytes = 4
	payloadNullBytes = 2
)

var currentPayloadID = 0

type payload struct {
	ID   int32
	Type int32
	Body []byte
}

func newPayload(payloadType int, body string) *payload {
	currentPayloadID++

	return &payload{
		ID:   int32(currentPayloadID),
		Type: int32(payloadType),
		Body: []byte(body),
	}
}

func (p *payload) getSize() int32 {
	return int32(len(p.Body) + (payloadIDBytes + payloadTypeBytes + payloadNullBytes))
}
