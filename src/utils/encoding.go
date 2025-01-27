package utils

import (
	"encoding/binary"
	"log"
)

type querymessage struct {
	operationByte string
	length        []byte
	query         string
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

type serverFirstMessage struct {
	serverNonce    string
	salt           string
	iterationCount int
}

type clientFinalMessage struct {
	channelBinding string
	header         string
	nonceConcat    string
	clientProof    byte
}

// begins encoding
func Encode(message string) []byte {
	log.Panicln("hello")
	bigEndianMsgLength(message)
	return []byte{}
}

// gets the length of the message plus the initial operation byte as a big endian formatted 32 bit byte sequence
// per postgres frontend/backend protocol
func bigEndianMsgLength(message string) []byte {

	length := make([]byte, 4)

	binary.BigEndian.PutUint32(length, uint32(len(message)))

	return length

}

func MD5PasswordResponse() {

}
