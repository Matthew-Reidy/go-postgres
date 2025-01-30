package utils

import (
	"crypto/tls"
	"encoding/binary"
	"log"

	pgtypes "github.com/Matthew-Reidy/go-postgres/src/types"
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

type ServerFirstMessage struct {
	ServerNonce    string
	Salt           string
	IterationCount int
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

var (
	globalConn *tls.Conn
)

// begins encoding
func Encode(message any, messageType string, conn *tls.Conn) []byte {
	globalConn = conn

	switch messageType {

	case "startup":

		signIn(message.(*pgtypes.Credentials))

	case "query":

	}

	return []byte{}
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

}

func signIn(message *pgtypes.Credentials) {
	start_up := []byte{}

	start_up = append(start_up, bigEndianMsgConverter(uint32(196608))...)

	start_up = append(start_up, []byte("user")...)

	start_up = append(start_up, 0x00)

	start_up = append(start_up, []byte(message.Database)...)

	start_up = append(start_up, 0x00, 0x00)

	start_up = append(bigEndianMsgConverter(uint32(len(start_up)+4)), start_up...)

	log.Println(start_up)

	b := make([]byte, 42)

	globalConn.Write(start_up)

	globalConn.Read(b)

	log.Println(string(b))

	authOptions := string(b)

	if authOptions == "R* SCRAM-SHA-256-PLUSSCRAM-SHA-256" {
		SCRAMRoutine()
	} else if authOptions == "MD5" {
		MD5PasswordRoutine()
	} else {
		//plaintext password auth...not implemented
		return
	}
}
