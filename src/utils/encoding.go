package utils

import (
	"encoding/binary"
	"log"
)

// begins encoding
func Encode() []byte {
	log.Panicln("hello")

	return []byte{}
}

// gets the length of the message plus the initial operation byte as a big endian formatted 32 bit byte sequence
// per postgres frontend/backend protocol
func bigEndianMsgLength(message string) []byte {

	length := make([]byte, 4)

	binary.BigEndian.PutUint32(length, uint32(len(message)+1))

	return length

}
