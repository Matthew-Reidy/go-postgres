package message

import (
	"encoding/binary"
	"fmt"

	pgtypes "github.com/Matthew-Reidy/go-postgres/types"
)

type querymessage struct {
	byte1  byte
	length uint32
	query  string
}

type startupmessage struct {
	length       int
	protocol_ver int
	username     string
	database     *string
	replication  *string
}

type clientFirstMessage struct {
	header      byte
	username    string
	clientNonce string
}

type clientFinalMessage struct {
	channelBinding string
	header         string
	nonceConcat    string
	clientProof    byte
}

type passwordMessage struct {
	byte1    byte
	length   uint32
	password string
}

const (
	AUTH_RESP = 'p'
)

// Encode builds the wire-format bytes for a message. It performs no I/O,
// so callers own the connection and are free to use it concurrently.
func Encode(message any, messageType any) ([]byte, error) {

	switch messageType {

	case pgtypes.STARTUP:

		creds, ok := message.(*pgtypes.Credentials)

		if !ok || creds == nil {
			return nil, fmt.Errorf("startup message requires *types.Credentials")
		}

		return startupMessage(creds), nil

	case pgtypes.QUERY:

		query, ok := message.(string)

		if !ok {
			return nil, fmt.Errorf("query message requires a string")
		}

		return queryMessage(query), nil

	}

	return nil, fmt.Errorf("unknown message type %q", messageType)
}

// gets the length of the message plus the initial operation byte as a big endian formatted 32 bit byte sequence
// per postgres frontend/backend protocol
func bigEndianMsgConverter(messageLen uint32) []byte {

	length := make([]byte, 4)

	binary.BigEndian.PutUint32(length, messageLen)

	return length

}

func MD5PasswordRoutine() {

}

func SCRAMRoutine() {
	fmt.Println("starting scram")
}

// startup message: int32 length, int32 protocol 3.0, then null terminated key/value pairs
func startupMessage(creds *pgtypes.Credentials) []byte {

	body := bigEndianMsgConverter(uint32(196608))

	for _, kv := range [][2]string{{"user", creds.Username}, {"database", creds.Database}} {
		body = append(body, kv[0]...)
		body = append(body, 0)
		body = append(body, kv[1]...)
		body = append(body, 0)
	}

	body = append(body, 0)

	return append(bigEndianMsgConverter(uint32(len(body)+4)), body...)
}

// simple query message: 'Q', int32 length, null terminated query string
func queryMessage(query string) []byte {

	msg := []byte{'Q'}

	msg = append(msg, bigEndianMsgConverter(uint32(len(query)+5))...)

	msg = append(msg, query...)

	return append(msg, 0)
}
