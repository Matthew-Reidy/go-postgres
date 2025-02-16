package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

type AuthenticationMessage struct {
	codeByte       rune
	allowedMethods []string
}

type ServerFirstMessage struct {
	ServerNonce    string
	Salt           string
	IterationCount int
}

func Decode(messageType string, respMessage []byte) interface{} {

	switch messageType {
	case "AuthOptions":

		return AuthOptions(respMessage)

	default:

		return nil

	}

}

func AuthOptions(respMessage []byte) *AuthenticationMessage {
	var opCode uint32

	buff := bytes.NewReader(respMessage)

	err := binary.Read(buff, binary.LittleEndian, &opCode)

	if err != nil {

		log.Fatalln("cant decode! ", err)

	}

	authOptions := bytes.Split(respMessage[9:len(respMessage)-1], []byte{0})

	strOpts := make([]string, len(authOptions))

	for i := range authOptions {

		strOpts[i] = string(authOptions[i])

	}

	authMessage := &AuthenticationMessage{
		codeByte:       rune(opCode),
		allowedMethods: strOpts,
	}

	return authMessage
}
