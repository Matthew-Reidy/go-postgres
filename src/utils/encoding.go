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
	length       []byte
	protocol_ver []byte
	username     string
	database     *string
	replication  *string
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
