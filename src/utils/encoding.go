package utils

import (
	"encoding/binary"
	"log"
	"net"
)

// begins encoding
func encode(conn *net.Conn) {
	log.Panicln("hello")
}

// gets the length of the message plus the initial operation byte as a big endian formatted 32 bit byte sequence
// per postgres frontend/backend protocol
func bigEndianMsgLength(message string) []byte {

	length := make([]byte, 4)

	binary.BigEndian.PutUint32(length, uint32(len(message)+1))

	return length

}
